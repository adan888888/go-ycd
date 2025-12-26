import axios from 'axios';

const instance = axios.create({
  baseURL: 'http://localhost:3000/api',
  timeout: 10000, // 10秒超时
});

instance.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = token;
  }
  return config;
});

// 响应拦截器：从原始响应文本中提取大整数，避免JSON解析时的精度丢失
instance.interceptors.response.use(response => {
  // 处理 /api/ycd/today/users 接口返回的数据
  if (response.config.url?.includes('/ycd/today/users') && response.data?.data) {
    if (Array.isArray(response.data.data)) {
      // 尝试从原始响应文本中提取正确的user_id
      // axios的response对象中，原始文本在response.request.responseText
      const originalText = (response as any).request?.responseText || (response as any).data;
      
      if (originalText && typeof originalText === 'string') {
        try {
          // 使用正则表达式提取所有user_id，避免JSON.parse时的精度丢失
          const userIdRegex = /"user_id"\s*:\s*(\d+)/g;
          const userIds: string[] = [];
          let match;
          while ((match = userIdRegex.exec(originalText)) !== null) {
            userIds.push(match[1]); // 提取原始字符串形式的user_id
          }
          
          // 将提取的user_id字符串替换到数据中
          if (userIds.length === response.data.data.length) {
            response.data.data = response.data.data.map((user: any, index: number) => ({
              ...user,
              user_id: userIds[index] // 使用从原始文本中提取的字符串
            }));
            return response;
          }
        } catch (e) {
          // 无法从原始文本提取，继续后续处理
        }
      }
      
      // 如果无法从原始文本提取，至少确保是字符串类型
      // 但此时精度可能已经丢失
      response.data.data = response.data.data.map((user: any) => {
        const userId = user.user_id;
        let userIdStr: string;
        if (typeof userId === 'number') {
          userIdStr = userId.toString();
        } else {
          userIdStr = String(userId);
        }
        return {
          ...user,
          user_id: userIdStr
        };
      });
    }
  }
  return response;
}, error => {
  // 统一错误处理
  return Promise.reject(error);
});

export default instance;
