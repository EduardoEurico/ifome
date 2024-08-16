<template>
  <div class="container">
    <div class="wrapper">
      <div class="card-switch">
        <label class="switch">
          <input type="checkbox" class="toggle" v-model="isSignUp">
          <span class="slider"></span>
          <span class="card-side"></span>
          <div class="flip-card__inner">
            <div class="flip-card__front">
              <div class="title">Log in</div>
              <form class="flip-card__form" @submit.prevent="handleLogin">
                <input class="flip-card__input" name="username" placeholder="Nome do Usuário" v-model="username" required>
                <input class="flip-card__input" name="password" placeholder="Senha" type="password" v-model="password" required>
                <button class="flip-card__btn">Let`s go!</button>
              </form>
              <router-link to="/auth/criarconta/restaurante" class="sidebar-link">Acesso de Restaurante</router-link>
            </div>
            <div class="flip-card__back">
              <div class="title">Sign up</div>
              <form class="flip-card__form" @submit.prevent="handleSignUp">
                <input class="flip-card__input" name="username" placeholder="Nome" v-model="username" required>
                <input class="flip-card__input" name="email" placeholder="E-mail" type="email" v-model="email" required>
                <input class="flip-card__input" name="password" placeholder="Senha" type="password" v-model="password" required>
                <input class="flip-card__input" name="confirmPassword" placeholder="Confirme a Senha" type="password" v-model="confirmPassword" required>
                <button class="flip-card__btn">Confirm!</button>
              </form>
            </div>
          </div>
        </label>
      </div>   
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      isSignUp: false // Para controlar se está na tela de login ou cadastro
    };
  },
  methods: {
    handleLogin() {
      const loginInfo = {
        username: this.username,
        password: this.password
      };

      axios.post('http://localhost:8080/login', loginInfo)
        .then(response => {
          console.log('Resposta da API:', response.data);
          const { token, userName } = response.data;

          if (token && userName) {
            localStorage.setItem('auth_token', token);
            localStorage.setItem('userName', userName);
            alert('Login realizado com sucesso!');
            this.$router.push({ name: 'Home' });
          } else {
            throw new Error('Dados de resposta da API inválidos');
          }
        })
        .catch(error => {
          console.error('Erro ao fazer login:', error);
          alert('Erro ao fazer login. Verifique suas credenciais.');
        });
    },
    handleSignUp() {
      if (this.password !== this.confirmPassword) {
        alert('As senhas não coincidem.');
        return;
      }

      const newUser = {
        username: this.username,
        email: this.email,
        password: this.password,
        confirmPassword: this.confirmPassword,
      };

      axios.post('http://localhost:8080/conta', newUser)
        .then(response => {
          console.log('Resposta da API:', response.data);
          alert('Conta criada com sucesso!');
        })
        .catch(error => {
          console.error('Erro ao criar a conta:', error.response ? error.response.data : error.message);
          alert('Erro ao criar a conta: ' + (error.response ? error.response.data.error : error.message));
        });
    }
  }
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');

* {
  font-family: 'Poppins', sans-serif; 
}
.container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 115vh;
  background: url('@/assets/fundo_criarconto.jpg') no-repeat center center fixed; /* Certifique-se que o caminho está correto */
  background-size: cover; /* Adicionado para cobrir a tela inteira */
}

.wrapper {
  --input-focus: #ff0000;
  --font-color: #323232;
  --font-color-sub: #000000;
  --bg-color: #fff;
  --bg-color-alt: #ffffff;
  --main-color: #ff0000;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.switch {
  transform: translateY(-200px);
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 30px;
  width: 50px;
  height: 20px;
}

.card-side::before {
  position: absolute;
  content: 'Log in';
  left: -70px;
  top: 0;
  width: 100px;
  text-decoration: underline;
  color:  #ffffff;
  font-weight: 600;
}

.card-side::after {
  position: absolute;
  content: 'Sign up';
  left: 70px;
  top: 0;
  width: 100px;
  text-decoration: none;
  color:  #ffffff;
  font-weight: 600;
}

.toggle {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  box-sizing: border-box;
  border-radius: 5px;
  border: 2px solid var(--main-color);
  box-shadow: 4px 4px var(--main-color);
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--bg-color);
  transition: 0.3s;
}

.slider:before {
  box-sizing: border-box;
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  border: 2px solid var(--main-color);
  border-radius: 5px;
  left: -2px;
  bottom: 2px;
  background-color: var(--bg-color);
  box-shadow: 0 3px 0 var(--main-color);
  transition: 0.3s;
}


.toggle:checked + .slider {
  background-color: var(--input-focus);
}

.toggle:checked + .slider:before {
  transform: translateX(30px);
}

.toggle:checked ~ .card-side:before {
  text-decoration: none;
}

.toggle:checked ~ .card-side:after {
  text-decoration: underline;
}

.flip-card__inner {
  width: 300px;
  height: 350px;
  position: relative;
  background-color: transparent;
  perspective: 1000px;
  text-align: center;
  transition: transform 0.8s;
  transform-style: preserve-3d;
}

.toggle:checked ~ .flip-card__inner {
  transform: rotateY(180deg);
}

.toggle:checked ~ .flip-card__front {
  box-shadow: none;
}

.flip-card__front, .flip-card__back {
  padding: 20px;
  position: absolute;
  display: flex;
  flex-direction: column;
  justify-content: center;
  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
  background: linear-gradient(180deg, #fcb045, #fd1d1d,#ff3737);
  gap: 20px;
  border-radius: 5px;
  border: 2px solid var(--main-color);
  box-shadow: 4px 4px var(--main-color);
}
.sidebar-link{
  color: #ffffff;
}

.flip-card__back {
  width: 100%;
  transform: rotateY(180deg);
}

.flip-card__form {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.title {
  margin: 20px 0 20px 0;
  font-size: 25px;
  font-weight: 900;
  text-align: center;
  background: #ff0000;
}

.flip-card__input {
  width: 250px;
  height: 40px;
  border-radius: 5px;
  border: 2px solid var(--main-color);
  background-color: var(--bg-color);
  box-shadow: 4px 4px var(--main-color);
  font-size: 15px;
  font-weight: 600;
  color: var(--font-color);
  padding: 5px 10px;
  outline: none;
}

.flip-card__input::placeholder {
  color: var(--font-color-sub);
  opacity: 0.8;
}

.flip-card__input:focus {
  border: 2px solid var(--input-focus);
}

.flip-card__btn:active, .button-confirm:active {
  box-shadow: 0px 0px var(--main-color);
  transform: translate(3px, 3px);
  background: #007bff;
}

.flip-card__btn {
  margin: 20px 0 20px 0;
  width: 120px;
  height: 40px;
  border-radius: 5px;
  border: 2px solid var(--main-color);
  background-color: var(--bg-color);
  box-shadow: 4px 4px var(--main-color);
  font-size: 17px;
  font-weight: 600;
  color: var(--font-color);
  cursor: pointer;
}
</style>
