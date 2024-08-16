import axios from 'axios';

// Criação de uma instância do Axios
const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api', // URL base da API
  headers: {
    'Content-Type': 'application/json'
  }
});

// Adiciona um interceptador de requisição
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('auth_token'); // Obtém o token do localStorage
    if (token) {
      config.headers.Authorization = `Bearer ${token}`; // Adiciona o token ao cabeçalho Authorization
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// Adiciona um interceptador de resposta
apiClient.interceptors.response.use(
  response => {
    // Supondo que o novo token venha na resposta bem-sucedida
    if (response.data.token) {
      localStorage.setItem('auth_token', response.data.token);
    }
    return response;
  },
  error => {
    if (error.response && error.response.status === 401) {
      const serverToken = error.response.data.token;
      const clientToken = localStorage.getItem('auth_token');

      if (serverToken !== clientToken) {
        localStorage.clear();
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;
