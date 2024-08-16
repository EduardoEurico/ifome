<template>
  <div class="container">
    <!-- Pop-up de boas-vindas ou de login/registro -->
    <div v-if="showPopup" class="login-popup">
      <p v-if="!userName">Você precisa estar logado para ver esta página.</p>
      <p v-else>
        Bem-vindo, {{ userName }}!
        <button class="close-button" @click="closePopup">X</button>
      </p>
      <div v-if="!userName">
        <!-- Adicionando classes para estilos específicos -->
        <router-link to="/auth/criarconta/restaurante" class="router-link-restaurante">Acesso de Restaurantes</router-link>
        <p>ou</p>
        <router-link to="/auth/criarconta" class="router-link-usuario">Acesso de Usuário</router-link>
      </div>
    </div>

    <!-- Componente Box renderizado aqui -->
    <Box />
  </div>
</template>

<script>
import axios from 'axios';
import Box from './boxes.vue'; // Importando o componente Box.vue

export default {
  name: 'HomePage',
  
  components: {
    Box // Declarando o componente Box
  },
  
  data() {
    return {
      userName: localStorage.getItem('userName') || '',
      showPopup: !localStorage.getItem('auth_token'),
      cancelTokenSource: axios.CancelToken.source(),
    };
  },
  methods: {
    async getUserID() {
      const token = localStorage.getItem('auth_token');
      if (!token) {
        this.showPopup = true;
        return null;
      }
      try {
        const jwtDecode = (await import('jwt-decode')).default; // Correção aqui
        const decodedToken = jwtDecode(token);
        return decodedToken.id;
      } catch (e) {
        console.error("Erro ao decodificar token:", e);
        this.showPopup = true;
        return null;
      }
    },
    fetchUserName() {
      this.getUserID().then(userID => {
        if (userID) {
          axios.get(`http://localhost:8080/user/${userID}`)
            .then(response => {
              this.userName = response.data.name;
              this.showPopup = false;
              localStorage.setItem('userName', response.data.name);
            })
            .catch(error => {
              console.error("Erro ao buscar o nome do usuário:", error);
              this.showPopup = true;
            });
        }
      });
    },
    closePopup() {
      this.showPopup = false;
    }
  },
  beforeUnmount() {
    this.cancelTokenSource.cancel('Operação cancelada pelo usuário.');
  },
  mounted() {
    this.getUserID().then(userID => {
      if (!userID) {
        this.showPopup = true;
      } else {
        this.fetchUserName();
      }
    });
  }
};
</script>

<style scoped>
/* Estilos específicos para esta página podem ser adicionados aqui */
.container{ 
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: url('@/assets/12-01.jpg') no-repeat center center fixed; /* Certifique-se que o caminho está correto */
  background-size: cover; /* Adicionado para cobrir a tela inteira */
}

.container::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 0, 0, 0.5); /* Vermelho com opacidade */
  z-index: -1;
}

.login-popup {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #323232;
  padding: 20px;
  border:solid #ff0000;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.close-button {
  font-size: 18px;
  color: #fff;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  margin: 0;
  position: absolute;
  top: 10px;
  right: 10px;
}

.close-button:hover {
  color: #555;
}

.router-link-restaurante {
  color: #ff0000; /* Cor cinza escuro (#323232) */
}

.router-link-usuario {
  color: #ff0000; /* Cor azul (#007bff) */
}

.router-link-restaurante:hover,
.router-link-usuario:hover {
  color: #a953db; /* Cor preta (#111) ao passar o mouse */
}
</style>