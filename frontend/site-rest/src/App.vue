<template>
  <div id="app">
    <!-- Sidebar -->
    <div class="sidebar-container" :class="{ 'open': isSidebarOpen }">
      <div class="sidebar">
        <router-link to="/home" class="sidebar-link">Home</router-link>
        <router-link to="/auth/criarconta" class="sidebar-link">Criar Conta/Login</router-link>
        <router-link :to="restauranteLink" class="sidebar-link">Página Restaurante</router-link>
        <button @click="logout" class="sidebar-link logout-link">Logout</button>
      </div>
    </div>
    <div id="main">
      <div class="sidebar-title">
        <h2 class="titulo">Ifome</h2>
      </div>
      <button class="openbtn" @click="toggleSidebar(!isSidebarOpen)">
        ☰ {{ isSidebarOpen ? 'Fechar' : 'Abrir' }}
      </button>
    </div>
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: 'App',
  computed: {
    restauranteLink() {
      try {
        // Obter o nome do usuário do localStorage
        const userName = localStorage.getItem('userName') ;
        // Construir o link para a página do restaurante
        return `/restaurante/${userName}`;
      } catch (error) {
        console.error('Erro ao acessar o localStorage:', error);
        return '/restaurante';
      }
    }
  },
  data() {
    return {
      isSidebarOpen: false
    };
  },
  methods: {
    logout() {
      try {
        // Limpa o localStorage
        localStorage.clear();

        // Redireciona o usuário para a página inicial
        this.$router.push('/home');
      } catch (error) {
        console.error('Erro ao limpar o localStorage:', error);
      }
    },
    toggleSidebar(state) {
      this.isSidebarOpen = state;
    }
  }
};
</script>

<style>
/* Estilos globais para toda a aplicação */
body {
  background-color: #fff; /* Cor preta para o fundo da página */
  color: #fff; /* Cor branca para o texto */
  font-family: Arial, sans-serif; /* Fonte padrão */
  justify-content: center;
  align-items: center;
  height: auto; /* Altura total da tela */
  margin: 0; /* Remove margens padrão */
}

/* Estilos para o título ao lado do botão de abrir sidebar */
.sidebar-title {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%; /* Ocupa toda a largura da tela */
  z-index: 1000; /* Garante que o título esteja acima do conteúdo principal */
  color: #ffffff; /* Cor azul para o título */
  background-color: #ff0000; /* Cor de fundo da barra */
  height: 50px; /* Altura igual à da sidebar */
  display: flex;
  align-items: center; /* Centraliza verticalmente */
  justify-content: center; /* Centraliza horizontalmente */
}

.sidebar-title h2 {
  margin: 0; /* Remove margem padrão do título */
}

.sidebar-title::before {
  content: ''; /* Conteúdo vazio */
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #ff0000; /* Cor da barra */
  z-index: -1; /* Coloca atrás do título */
}

/* Estilo para os links da sidebar */
.sidebar a.sidebar-link,
.sidebar button.sidebar-link {
  padding: 10px 15px;
  text-decoration: none;
  font-size: 18px;
  color: #000000; /* Cor branca para os links normais */
  display: block;
  transition: 0.3s;
}

.sidebar a.sidebar-link:hover,
.sidebar button.sidebar-link:hover {
  color: #ff0000; /* Cor azul ao passar o mouse */
}

/* Estilo específico para o link de logout */
.sidebar a.logout-link {
  color: #fff; /* Cor branca para o link de logout */
}

/* Estilos do sidebar e do botão */
.sidebar-container {
  position: fixed;
  top: 0;
  left: -250px;
  height: 100vh;
  width: 250px;
  background-color: #fff; /* Cor preta para o sidebar */
  transition: left 0.5s ease;
  z-index: 1000; /* Garante que o sidebar esteja acima do conteúdo principal */
}

.sidebar-container.open {
  left: 0;
}

.sidebar {
  padding-top: 60px;
}

.openbtn {
  font-size: 20px;
  width: 250px;
  cursor: pointer;
  background-color: #ff0000; /* Cor preta para o botão */
  color: #fff; /* Cor branca para o texto do botão */
  border: none;
  padding: 10px 15px;
  border-radius: 1px;
  margin: 0;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1000; /* Garante que o botão esteja acima do conteúdo principal */
}

.openbtn:hover {
  background-color: #9c0000; /* Cor azul ao passar o mouse */
}
</style>
