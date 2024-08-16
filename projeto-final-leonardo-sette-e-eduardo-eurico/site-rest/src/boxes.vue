<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div>
    <div class="card-container">
      <div v-for="restaurante in restaurante" :key="restaurante.restname">
        <div class="card" @click="goToRestaurante(restaurante.restname)">
          <div class="card-image">
            <!-- Adicione a imagem aqui, verificando se existe -->
            <img v-if="restaurante.image" :src="restaurante.image" alt="Imagem do Restaurante">
          </div>
          <div class="card-content">
            <h2 class="card-title">{{ restaurante.restname }}</h2>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      restaurante: []
    };
  },
  created() {
    this.fetchRestaurantes();
  },
  methods: {
    async fetchRestaurantes() {
      try {
        const response = await axios.get('http://localhost:8080/restaurante');
        this.restaurante = response.data;
      } catch (error) {
        console.error('Error fetching restaurante:', error);
      }
    },
    goToRestaurante(restname) {
      this.$router.push(`/restaurante/${restname}`);
    }
  }
};
</script>

<style scoped>
.card-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center; /* Centraliza as boxes horizontalmente */
  align-items: flex-start; /* Alinha as boxes no topo verticalmente */
  gap: 30px; /* Espaçamento entre as boxes */
  margin-top: 200px; /* Espaçamento acima do contêiner */
  padding: 20px; /* Espaçamento interno do contêiner */
}

.card {
  width: 190px;
  height: 154px;
  border-radius: 30px;
  background: #ffffff;
  box-shadow: 5px 5px 10px rgb(255, 242, 67), -15px -15px 30px rgb(255, 208, 0);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  position: relative;
  cursor: pointer;
}

.card-image {
  width: 100%; /* Ajuste para a largura desejada ou use 100% para largura total */
  height: auto; /* Altura automática para manter a proporção */
  border-radius: 30px 30px 0 0;
  overflow: hidden; /* Garante que a imagem não ultrapasse o border-radius */
}

.card-image img {
  width: 100%; /* Faz a imagem se ajustar à largura do contêiner */
  height: auto; /* Mantém a proporção da imagem */
  display: block; /* Remove qualquer espaço abaixo da imagem */
}

.card-content {
  padding: 10px;
  color: #ffffff;
  text-align: center;
  background: rgb(255, 0, 0);
}

.card-title {
  margin: 0;
  font-size: 18px;
}

.card-rating {
  margin-top: 5px;
}

.card-rating span {
  font-size: 20px;
}
</style>
