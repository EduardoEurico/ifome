<template>
    <div>
      <div class="sidebar-container" :class="{ open: isSidebarOpen }">
        <div class="sidebar">
          <router-link to="/meu-perfil" v-if="isSidebarOpen">Meu Perfil</router-link>
          <router-link to="/restaurantes" v-if="isSidebarOpen">Restaurantes</router-link>
        </div>
      </div>
      <div id="main">
        <button class="openbtn" @click="toggleSidebar(!isSidebarOpen)">
          ☰ {{ isSidebarOpen ? 'Fechar' : 'Abrir' }}
        </button>
        <div class="content-container">
          <div class="content">
            <h2>Bem-vindo à aplicação!</h2>
            <p>Conteúdo principal da página.</p>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        isSidebarOpen: false
      }
    },
    methods: {
      toggleSidebar(isOpen) {
        this.isSidebarOpen = isOpen
      }
    }
  }
  </script>
  
  <style scoped>
  /* Mesmos estilos do arquivo fornecido */
  .sidebar-container {
    position: fixed;
    top: 0;
    left: -250px;
    height: 100vh;
    transition: left 0.5s ease;
  }
  
  .sidebar-container.open {
    left: 0;
  }
  
  .sidebar {
    background-color: #111;
    padding-top: 60px;
    width: 250px;
    height: 100%;
    position: relative;
  }
  
  .sidebar a {
    padding: 8px 8px 8px 32px;
    text-decoration: none;
    font-size: 25px;
    color: #818181;
    display: block;
    transition: 0.3s;
  }
  
  .sidebar a:hover {
    color: #f1f1f1;
  }
  
  .openbtn {
    font-size: 20px;
    width: 250px;
    cursor: pointer;
    background-color: #111;
    color: white;
    border: none;
    padding: 10px 15px;
    border-radius: 1px;
    margin: 0;
    position: fixed;
    top: 0;
    left: 0;
  }
  
  .openbtn:hover {
    background-color: #444;
  }
  
  .content-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: calc(100vh - 40px);
  }
  
  .content {
    text-align: center;
  }
  
  #main {
    padding: 20px;
    margin-left: 0;
    background-color: black;
    color: white;
    height: 100vh;
    transition: margin-left 0.5s;
  }
  </style>
  