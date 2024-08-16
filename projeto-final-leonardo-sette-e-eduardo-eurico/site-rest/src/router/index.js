// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import Home from '../HomePage.vue';
import CriarConta from '../CriarConta.vue';
import CriarContaRest from '../CriarContaRest.vue';
import Usuario from '../UserPage.vue';
import Restaurante from '../RestaurantPage.vue';

const routes = [
  {
    path: '/home',
    name: 'Home',
    component: Home
  },
  {
    path: '/',
    redirect: '/home',
  },
  {
    path: '/auth/criarconta',
    name: 'CriarConta',
    component: CriarConta
  },
  {
    path: '/auth/criarconta/restaurante',
    name: 'CriarContaRest',
    component: CriarContaRest
  },
  {
    path: '/usuario',
    name: 'Usuario',
    component: Usuario,
    meta: { requiresAuth: true } // Marca a rota como protegida
  },
  {
    path: '/restaurante/:userName',
    name: 'Restaurante',
    component: Restaurante,
    meta: { requiresAuth: true }, // Marca a rota como protegida
    props: true 
  }
];

// Função para verificar se o usuário está autenticado
function isAuthenticated() {
  // Supondo que o token de autenticação seja armazenado no localStorage
  return !!localStorage.getItem('auth_token');
}

// Adicionar navegação guard para verificar a autenticação
const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // Verificar se o usuário está autenticado
    if (!isAuthenticated()) {
      next('/auth/criarconta'); // Redireciona para a página de login se não estiver autenticado
    } else {
      next(); // Continue para a rota solicitada
    }
  } else {
    next(); // Continue para a rota solicitada
  }
});

export default router;
