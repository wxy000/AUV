// 公共配置
const dayjs = window.dayjs;
const utcPlugin = window.dayjs_plugin_utc;
dayjs.extend(utcPlugin);

dayjs.locale('zh-cn');

// 全局配置
const API_BASE = '/api/v1';

// 通用请求头
function getDefaultHeaders(includeAuth = true) {
  const headers = {
    'Content-Type': 'application/json'
  };
  if (includeAuth) {
    headers.Authorization = 'Bearer ' + (localStorage.getItem('jwt') || '');
  }
  return headers;
}

async function request(endpoint, method = 'GET', body = null, includeAuth = true) {
  try {
    const response = await fetch(`${API_BASE}${endpoint}`, {
      method,
      headers: getDefaultHeaders(includeAuth),
      body: body ? JSON.stringify(body) : null
    });
    console.log(`API请求: ${endpoint}, 状态: ${response.status}`);

    if (response.status === 401) {
      await handleAuthError();
    }

    if (response.status === 429) {
      Swal.fire({
        icon: 'warning',
        title: '操作提示',
        text: '请求过于频繁，请稍后再试',
        confirmButtonText: '知道了'
      });
    }

    if (!response.ok) {
      throw {status: response.status,message: await response.text(),handled: true}
    }

    const result = await response.json();
    
    // 统一校验响应格式
    if (typeof result.code === 'undefined' && !result.message && !result.data) {
      throw { 
        code: 500, 
        message: '非标准接口响应格式',
        handled: false
      };
    }
    
    return result.data;
  } catch (error) {
    console.error(`API请求失败: ${error}`);
    if (!error.handled) {
      Swal.fire({
        icon: error.code === 401 ? 'warning' : 'error',
        title: `操作失败（${error.code}）`,
        text: error.message || '请求处理失败',
        confirmButtonText: '确定'
      });
    }
    throw error;
  }
}

// 用户相关接口
async function getUser(userId) {
  return request(`/getUser/${userId}`);
}

async function loadUsers(showActive = true) {
  const endpoint = `/${showActive ? 'getUsers' : 'getInactiveUsers'}`;
  return request(endpoint);
}

async function saveUser(userData, userId = '') {
  const endpoint = `/${userId ? `updateUser/${userId}` : 'createUser'}`;
  return request(endpoint, 'POST', userData);
}

// 认证相关
function handleAuthError(redirect = true) {
  localStorage.removeItem('jwt');
  localStorage.removeItem('userId');
  if(redirect) {
    Swal.fire({
      icon: 'error',
      title: '会话过期，请重新登录',
      showConfirmButton: false,
      timer: 1200
    });
    window.location.href = '../admin/login.html';
  }
}

// 初始化检查登录状态
function checkLogin() {
  const jwt = localStorage.getItem('jwt');
  const userId = localStorage.getItem('userId');
  
  if(!jwt || !userId) {
    window.location.href = '../admin/login.html';
    return;
  }

  // 直接通过接口获取用户信息
  getUser(userId)
    .then(user => {
      document.getElementById('currentUser').textContent = user.realname || user.username || '未知用户';
    });
}

// 登录方法
async function loginUser(username, password) {
  try {
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: getDefaultHeaders(false),
      body: JSON.stringify({ username, password })
    });
    
    const data = await response.json();

    if (!response.ok) {
      // 登录失败专用错误格式
      throw {status: response.status===401,message: response.text() || '认证失败',handled: true}
    }

    localStorage.setItem('jwt', data.data.access_token);
    localStorage.setItem('userId', data.data.user_id);
    return data.data;
  } catch (error) {
    console.error('登录错误:', error);
    
    if (error.status === 429) {
      Swal.fire({
        icon: 'error',
        title: '登录失败',
        text: '请求过于频繁，请稍后再试',
        confirmButtonText: '知道了'
      });
      return;
    }

    // 仅处理非认证头的401错误
    if (error.isAuthError) {
      Swal.fire({
        icon: 'error',
        title: '登录失败',
        text: error.message || '用户名密码错误',
        confirmButtonText: '重新输入'
      });
    }
    
    Swal.fire({
      icon: 'error',
      title: '登录失败',
      text: error.message || '请检查用户名和密码',
      confirmButtonText: '重新输入'
    });

    throw error;
  }
}

async function logout() {
    localStorage.removeItem('jwt');
    localStorage.removeItem('userId');
    window.location.href = '../admin/login.html';
}

async function deleteUser(userId) {
    const { isConfirmed } = await Swal.fire({
        icon: 'warning',
        title: '确认删除',
        text: '确定要永久删除该用户吗？此操作不可撤销！',
        showCancelButton: true,
        confirmButtonText: '确认删除',
        cancelButtonText: '取消操作',
        confirmButtonColor: '#d33',
        cancelButtonColor: '#3085d6'
    });
    if (!isConfirmed) return;
    try {
        await request(`/deleteUser/${userId}`, 'DELETE');
        Swal.fire({
          icon: 'success',
          title: '删除成功',
          showConfirmButton: false,
          timer: 1500
        });
        location.reload();
    } catch (error) {
        Swal.fire({
          icon: 'error',
          title: '删除失败',
          text: error.message || '删除操作失败，请稍后重试'
        });
    }
}

async function updateUserStatus(userId, isActive) {
    try {
        await request(`/updateUserStatus/${userId}`, 'POST', { isActive });
        Swal.fire({
          icon: 'success',
          title: '状态已更新',
          showConfirmButton: false,
          timer: 1200
        });
        location.reload();
    } catch (error) {
        Swal.fire({
          icon: 'error',
          title: '操作失败',
          text: error.message || '状态更新出现问题，请稍后重试'
        });
    }
}

// 暴露公共方法
window.CommonAPI = {
    getUser,
    loadUsers,
    saveUser,
    checkLogin,
    handleAuthError,
    loginUser,
    logout,
    deleteUser,
    updateUserStatus
  };