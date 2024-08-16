<template>
  <div class="container">
    <header class="header">
      <div class="image-upload" @click="triggerFileInput" @dragover.prevent @drop="onFileDrop" v-if="isOwner">
        <input type="file" id="image" @change="handleFileChange" ref="fileInput" />
        <span>Clique ou arraste uma imagem aqui</span>
      </div>
      <div v-if="imagemBase64">
          <img :src="imagemBase64" class="preview-image"/> <!-- Remover 'imagemBase64' + -->
        </div>
        <div v-else>
          
        </div>
      <h1 class="logo"> {{ ownerUserName }}</h1>
      <nav class="nav" v-if="isOwner">
        <button class="nav-button" @click="adicionarLista">Adicionar Lista</button>
        <button class="nav-button" @click="saveLists">Salvar Listas</button>
      </nav>
    </header>

    <main class="main">
      <div>
        <div class="listas-wrapper">
          <div v-for="(lista, index) in lista" :key="index" class="lista">
            <div class="lista-header">
              <h3 class="lista-nome">{{ lista.nome }}</h3>
              <div v-if="isOwner">
                <button @click="editarLista(index)" class="editar-lista"><i class="fas fa-pencil-alt"></i></button>
                <button @click="adicionarProduto(index)" class="adicionar-produto"><i class="fas fa-plus"></i></button>
                <button @click="removerLista(index)" class="trash-button"><i class="fas fa-trash"></i></button>
              </div>
            </div>
            <ul class="lista-itens">
              <li v-for="item in lista.itens" :key="item.nome">
                <strong>{{ item.nome }}</strong><br>
                Valor: {{ item.valor }}<br>
                <div class="descricao-item">Descrição: {{ item.descricao }}</div><br>
                <img :src="item.image" alt="Foto do produto" v-if="item.image"/>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </main>

    <div v-if="feedbackMessage" :class="`feedback ${feedbackType}`">
      {{ feedbackMessage }}
    </div>
  </div>
</template>





<script>
import axios from 'axios';

function makeAuthenticatedRequest(endpoint) {
    const token = localStorage.getItem('auth_token');
    return axios.get(endpoint, {
        headers: {
            'Authorization': ` ${token}`
        }
    });
}

// Exemplo de chamada
makeAuthenticatedRequest('/auth/checkPermissions')
    .then(response => {
        console.log("Resposta da API:", response.data);
    })
    .catch(error => {
        console.error("Erro na API:", error);
    });

export default {
  data() {
    return {
      ownerUserName: '',
      authToken: localStorage.getItem('auth_token') || '',
      userName: localStorage.getItem('userName') || '',
      isOwner: false,
      lista: [],
      feedbackMessage: '',
      feedbackType: '',
      listaCategoria: '',
      image: '',
      imagemBase64: ''
    };
  },

  methods: {
    async verifyOwner() {
    try {
      const urlUserName = this.$route.params.userName;
      const dataOwner = { urlUserName };

      const response = await axios.post('http://localhost:8080/filterUrl', dataOwner, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': ` ${this.authToken}`
        }
      });

      console.log("Resposta da API (verifyOwner):", response.data);

      if (response.data.IsOwner) {
        this.isOwner = true;
      } else {
        this.isOwner = false;
      }
    } catch (error) {
      console.error('Erro ao verificar proprietário:', error);
    }
  },
    async verifyAuthToken() {
      try {
        const urlUserName = this.$route.params.userName;
        console.log("Valor de urlUserName:", urlUserName);

        if (this.authToken) {
          const response = await axios.get('http://localhost:8080/isOwner', {
            headers: { 'Authorization': ` ${this.authToken}` },
            params: { urlUserName: urlUserName }
          });
          console.log("Resposta da API (verifyAuthToken):", response.data);

          if (response.data.IsOwner) {
            this.isOwner = true;
            this.ownerUserName = urlUserName;
          } else {
            this.isOwner = false;
            this.ownerUserName = urlUserName;
          }
        } else {
          this.isOwner = false;
          this.ownerUserName = urlUserName;
        }
      } catch (error) {
        console.error('Erro ao verificar token (verifyAuthToken):', error);
        console.log("Token usado:", this.authToken);
        console.log("UrlUserName usado:", this.$route.params.userName);
      }
    },

    async loadLists() {
      const ownerUserName = this.$route.params.userName;

      if (ownerUserName) {
        try {
          const response = await axios.get(`http://localhost:8080/restaurante/${ownerUserName}/lista/get`, {
            headers: { 'Authorization': ` ${this.authToken}` }
          });

          if (response.data && response.data.listas) {
            this.lista = response.data.listas;
          } else {
            console.error('A resposta não contém a lista esperada:', response.data);
          }

          const stringResponse = await axios.get(`http://localhost:8080/restaurante/${ownerUserName}/pegarimagem`, {
            headers: { 'Authorization': ` ${this.authToken}` }
          });
          console.log('Resposta completa da imagem:', stringResponse);

          if (stringResponse.data && stringResponse.data.image) {
            const imageString = stringResponse.data.image.trim();
            if (imageString !== '') {
              this.imagemBase64 = imageString;
            } else {
              console.error('A string da imagem está vazia após trim:', stringResponse.data);
            }
          } else {
            console.error('A resposta não contém a chave "image" ou está vazia:', stringResponse.data);
          }
        } catch (error) {
          console.error('Erro ao carregar listas ou imagem:', error);
          if (error.response) {
            console.error('Status:', error.response.status);
            console.error('Dados:', error.response.data);
          }
        }
      }
    },

    adicionarLista() {
      if (this.isOwner) {
        const nome = prompt('Digite o nome da nova lista:');

        if (nome) {
          this.lista.push({ nome, itens: [] });
          this.saveLists();
        }
      }
    },

    async saveLists() {
      if (this.isOwner) {
        try {
          for (let i = 0; i < this.lista.length; i++) {
            const lista = this.lista[i];
            if (!lista.itens) {
              lista.itens = [];
            }

            const data = {
              nome: lista.nome,
              categoria: lista.categoria,
              itens: lista.itens.map(item => ({
                nome: item.nome,
                valor: item.valor,
                descricao: item.descricao,
              }))
            };
            console.log('Dados a serem enviados:', data);
            const response = await axios.post(`http://localhost:8080/saveLists/${this.ownerUserName}`, data, {
              headers: { 'Authorization': ` ${this.authToken}` }
            });

            if (response.status === 200) {
              if (response.data.message === "List added successfully") {
                this.feedbackMessage = 'Lista adicionada com sucesso!';
              } else if (response.data.message === "List updated successfully") {
                this.feedbackMessage = 'Lista atualizada com sucesso!';
              }
              this.feedbackType = 'success';
            } else {
              this.feedbackMessage = 'Erro ao salvar listas.';
              this.feedbackType = 'error';
            }
          }
        } catch (error) {
          console.error('Erro ao salvar listas:', error);
          this.feedbackMessage = 'Erro ao salvar listas.';
          this.feedbackType = 'error';
        } finally {
          setTimeout(() => {
            this.feedbackMessage = '';
            this.feedbackType = '';
          }, 5000);
        }
      }
    },

    async adicionarProduto(index) {
      if (this.isOwner) {
        const nomeProduto = prompt('Digite o nome do produto:');
        const valorProduto = parseFloat(prompt('Digite o valor do produto:'));
        const descricaoProduto = prompt('Digite a descrição do produto:');
        if (nomeProduto && valorProduto && descricaoProduto) {
          this.lista[index].itens.push({
            nome: nomeProduto,
            valor: valorProduto,
            descricao: descricaoProduto
          });

          this.saveLists();
        }
      }
    },

    editarLista(index) {
      if (this.isOwner) {
        const nome = prompt('Digite o novo nome da lista:', this.lista[index].nome);

        if (nome) {
          this.lista[index].nome = nome;
          this.saveLists();
        }
      }
    },

    async removerLista(index) {
      if (this.isOwner && this.lista) {
        try {
          const nomeLista = this.lista[index].nome;
          const response = await axios.delete(`http://localhost:8080/restaurante/${this.ownerUserName}/lista/${nomeLista}`, {
            headers: { 'Authorization': ` ${this.authToken}` }
          });
          if (response.status === 200) {
            this.feedbackMessage = 'Lista removida com sucesso!';
            this.feedbackType = 'success';
          } else {
            this.feedbackMessage = 'Erro ao remover lista.';
            this.feedbackType = 'error';
          }
          this.loadLists();
        } catch (error) {
          console.error('Erro ao remover lista:', error);
        } finally {
          setTimeout(() => {
            this.feedbackMessage = '';
            this.feedbackType = '';
          }, 5000);
        }
      }
    },

    convertBase64ToBlob(base64Image) {
      const parts = base64Image.split(';base64,');
      const contentType = parts[0].split(':')[1];
      const raw = window.atob(parts[1]);
      const rawLength = raw.length;
      let uInt8Array = new Uint8Array(rawLength);

      for (let i = 0; i < rawLength; ++i) {
        uInt8Array[i] = raw.charCodeAt(i);
      }

      return new Blob([uInt8Array], { type: contentType });
    },

    handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.image = e.target.result;
          this.uploadImage();
        };
        reader.readAsDataURL(file);
      }
    },

    uploadImage() {
      if (!this.image) {
        console.error("Nenhuma imagem selecionada para upload.");
        return;
      }

      const imageData = {
        image: this.image
      };

      axios.post(`http://localhost:8080/restaurante/${this.ownerUserName}/imagem`, imageData, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': ` ${this.authToken}`
        }
      }).then(response => {
        console.log("Upload bem-sucedido", response);
      }).catch(error => {
        console.error("Erro no upload", error);
      });
    },

    triggerFileInput() {
      this.$refs.fileInput.click();
    },

    onFileChange(e) {
      const file = e.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.image = e.target.result;
        };
        reader.readAsDataURL(file);
      }
    },

    onFileDrop(e) {
      e.preventDefault();
      const file = e.dataTransfer.files[0];
      this.onFileChange({ target: { files: [file] } });
    }
  },

  mounted() {
    this.verifyAuthToken();
    this.loadLists();
  }
};
</script>





<style scoped>
/* Define a fonte e a cor do texto para o conteúdo dentro do container */
.container {
  font-family: Arial, sans-serif;
  color: #333;
  /* Outros estilos necessários */
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: url('@/assets/mesa.jpg') no-repeat center center fixed; /* Certifique-se que o caminho está correto */
  background-size: cover; /* Adicionado para cobrir a tela inteira */
}


/* Estiliza o cabeçalho fixo */
.header {
  position: fixed; /* Fixa o cabeçalho na posição */
  top: 50px; /* Define a distância do topo da página */
  left: 0; /* Alinha o cabeçalho à esquerda */
  width: 100%; /* Define a largura do cabeçalho */
  display: flex; /* Permite o uso de flexbox para alinhamento dos itens */
  align-items: center; /* Centraliza os itens verticalmente */
  flex-wrap: wrap; /* Permite que os itens quebrem para a próxima linha em telas menores */
  padding: 10px 1px; /* Espaçamento interno do cabeçalho */
  box-shadow: 0 2px 4px rgba(0,0,0,0.1); /* Adiciona uma sombra sutil */
  height: 30vh;
  background: url('@/assets/joinha.jpg') no-repeat center center fixed; /* Certifique-se que o caminho está correto */
  background-size: cover; /* Adicionado para cobrir a tela inteira */}

/* Estiliza o componente de upload de imagem */
.image-upload {
  position: relative; /* Define uma posição relativa para o componente */
  width: 150px; /* Largura do componente */
  height: 150px; /* Altura do componente */
  border: 2px dashed #ccc; /* Estilo da borda */
  display: flex; /* Permite o uso de flexbox para alinhamento */
  align-items: center; /* Centraliza os itens verticalmente */
  justify-content: center; /* Centraliza os itens horizontalmente */
  cursor: pointer; /* Transforma o cursor em um ponteiro */
  overflow: hidden; /* Esconde qualquer conteúdo que exceda o limite */
  margin-left: 3%; /* Move os botões para o limite direito */

}

/* Esconde o input padrão do componente de upload de imagem */
.image-upload input {

display: none;
}

/* Estiliza a prévia da imagem */
.image-preview {
  position: relative; /* Define uma posição relativa para o componente */
  width: 100%; /* Largura da prévia da imagem */
  height: 100%; /* Altura da prévia da imagem */
  overflow: hidden; /* Esconde qualquer conteúdo que excede o limite */
}

/* Estiliza a imagem na prévia */
.preview-image {
  position: relative; /* Define uma posição relativa para o componente */
  width: 150px; /* Largura do componente */
  height: 150px; /* Altura do componente */
  border: 2px dashed #ccc; /* Estilo da borda */
  display: flex; /* Permite o uso de flexbox para alinhamento */
  align-items: center; /* Centraliza os itens verticalmente */
  justify-content: center; /* Centraliza os itens horizontalmente */
  cursor: pointer; /* Transforma o cursor em um ponteiro */
  overflow: hidden; /* Esconde qualquer conteúdo que exceda o limite */
  margin-left: 15%; /* Move os botões para o limite direito */
    
}
/* Estiliza o logo */
.logo {
font-size: 40px; /* Tamanho da fonte */
font-weight: bold; /* Peso da fonte */
color: #ffffff; /* Cor do texto */
margin-bottom: 25px; /* Espaçamento abaixo do logo */
white-space: nowrap; /* Impede a quebra de linha do texto */
display: flex; /* Permite o uso de flexbox para alinhamento */
margin-left: 3%; /* Move os botões para o limite direito */

}

/* Estiliza a navegação */
.nav {
  display: flex; /* Permite o uso de flexbox para alinhamento */
  gap: 10px; /* Espaço entre os itens da navegação */
  margin-left: 30%; /* Move os botões para o limite direito */
}

/* Estiliza os botões da navegação */
.nav-button {
  background-color: #ff0000; /* Cor de fundo do botão */
  color: #fff; /* Cor do texto do botão */
  border: none; /* Remove a borda */
  padding: 10px 20px; /* Espaçamento interno do botão */
  border-radius: 5px; /* Bordas arredondadas */
  cursor: pointer; /* Transforma o cursor em um ponteiro */
}

/* Estiliza o efeito hover nos botões da navegação */
.nav-button:hover {
  background-color: #cc0000; /* Cor de fundo alterada no hover */
}

/* Estiliza o container principal */
.main {
  margin-top: 250px; /* Espaçamento do topo */
}

/* Estiliza cada lista */
.lista {
  background-color: #fff; /* Cor de fundo da lista */
  border: 1px solid #ddd; /* Estilo da borda */
  border-radius: 10px; /* Bordas arredondadas */
  box-shadow: 0 4px 8px rgba(0,0,0,0.1); /* Sombra sutil */
  margin-bottom: 20px; /* Espaço abaixo de cada lista */
}
/* Estiliza o cabeçalho da lista */
.lista-header {
  display: flex; /* Utiliza flexbox para alinhar os itens */
  align-items: center; /* Centraliza os itens verticalmente */
  justify-content: space-between; /* Distribui os itens horizontalmente */
  margin-bottom: 10px; /* Espaço abaixo do cabeçalho */
  margin-left: 10px; /* Espaçamento à esquerda */
}

.lista-nome{
  padding-right: 100px; /* Aumenta o espaço à direita */
}

.listas-wrapper {
  display: flex; /* Exibe as listas lado a lado */
  flex-wrap: wrap; /* Permite que as listas quebrem para a próxima linha */
  gap: 20px; /* Espaçamento entre as listas */
}
.lista-header button {
  margin-top: 5px; /* Espaço abaixo do cabeçalho */
  margin-left: 5px; /* Espaçamento à esquerda dos botões */
  margin-right: auto;
}
/* Estiliza o título da lista */
.lista h3 {
  margin: 0; /* Remove as margens padrão */
  font-size: 18px; /* Tamanho da fonte */
  color: #333; /* Cor do texto */
}

/* Estiliza a lista de itens */
.lista-itens {
  list-style: none; /* Remove os marcadores da lista */
  padding: 0; /* Remove o espaçamento interno */
  margin-top: 10px; /* Espaço acima dos itens */
}

/* Estiliza cada item da lista */
.lista-itens li {
  font-family: Arial, sans-serif;
  color: #333;
  background-color: #f0f0f0; /* Cor de fundo do item */
  margin-bottom: 10px; /* Espaço abaixo de cada item */
  padding: 10px; /* Espaçamento interno */
  border-radius: 5px; /* Bordas arredondadas */
  max-width: 280px; /* Largura máxima da lista */
  overflow: hidden; /* Esconde o conteúdo que excede a largura máxima */
  word-wrap: break-word; /* Quebra de palavra quando necessário */
}
.lista-itens li .descricao-item {
  white-space: pre-wrap; /* Preserva quebras de linha no texto */
}

/* Estiliza o botão de adicionar produto */
.editar-lista {
  background-color: #007bff; /* Cor de fundo do botão */
  color: #fff; /* Cor do texto */
  border: none; /* Remove a borda */
  padding: 10px; /* Espaçamento interno */
  border-radius: 5px; /* Bordas arredondadas */
  cursor: pointer; /* Transforma o cursor em um ponteiro */
  margin-right: 10px; /* Espaçamento à direita do botão */
}

/* Estiliza o ícone do botão de adicionar produto */
.adicionar-produto {
  background-color: #28a745; /* Cor de fundo do botão */
  color: #fff; /* Cor do texto */
  border: none; /* Remove a borda */
  padding: 10px; /* Espaçamento interno */
  border-radius: 5px; /* Bordas arredondadas */
  cursor: pointer; /* Transforma o cursor em um ponteiro */
  margin-right: 10px; /* Espaçamento à direita do botão */
}
.descricao-item {
  max-width: 100%; /* Largura máxima da descrição */
  max-height: 100px; /* Altura máxima da descrição */
  overflow: hidden; /* Esconde o conteúdo que excede a altura máxima */
  overflow-wrap: break-word; /* Quebra de palavra quando necessário */
  word-wrap: break-word; /* Quebra de palavra para navegadores antigos */
  display: inline-block; /* Força o elemento a se ajustar ao conteúdo */
  
}
.descricao-item.truncada {
  white-space: normal; /* Permite que o texto quebre para a próxima linha */
}

/* Estiliza o ícone do botão de excluir */
.trash-button {
  background-color: #dc3545; /* Cor de fundo do botão */
  color: #fff; /* Cor do texto */
  border: none; /* Remove a borda */
  padding: 10px; /* Espaçamento interno */
  border-radius: 5px; /* Bordas arredondadas */
  cursor: pointer; /* Transforma o cursor em um ponteiro */
}


/* Estiliza o efeito hover no botão de adicionar produto */
.adicionar-produto:hover {
  background-color: #0056b3; /* Cor de fundo alterada no hover */
}

/* Estiliza o ícone do botão de excluir */


/* Estiliza o ícone dentro do botão de excluir */
.trash-button i,
.editar-lista i,
.adicionar-produto i {
  font-size: 18px; /* Tamanho do ícone */
}

/* Estiliza a mensagem para usuários sem permissão */
.no-permission {
  font-size: 16px; /* Tamanho da fonte */
  color: #666; /* Cor do texto */
}

/* Estiliza o feedback exibido na tela */
.feedback {
  padding: 10px; /* Espaçamento interno */
  margin-top: 20px; /* Espaço acima do feedback */
  border-radius: 5px; /* Bordas arredondadas */
  text-align: center; /* Alinhamento central do texto */
}

/* Estiliza o feedback de sucesso */
.success {
  background-color: #d4edda; /* Cor de fundo */
  color: #155724; /* Cor do texto */
}

/* Estiliza o feedback de erro */
.error {
  background-color: #f8d7da; /* Cor de fundo */
  color: #721c24; /* Cor do texto */
}

/* Estiliza o feedback de informação */
.info {
  background-color: #cce5ff; /* Cor de fundo */
  color: #004085; /* Cor do texto */
}

@media (max-width: 768px) {
.lista {
  width: 100%; /* Em telas menores, ocupa 100% da largura */
  max-width: 100%; /* Largura máxima das listas */
}

.nav {
  flex-direction: column; /* Navegação em coluna em telas menores */
}

.lista {
  max-width: 50%; /* Lista ocupa 100% da largura em telas menores */
}
}
</style>