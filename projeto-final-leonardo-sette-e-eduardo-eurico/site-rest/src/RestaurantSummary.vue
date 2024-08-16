<template>
    <div class="product-grid">
      <div class="product-card" v-for="n in 3" :key="n">
        <div class="image-placeholder">
          <span>450 x 300</span>
        </div>
        <div class="sale-badge">Restaurante</div>
        <div class="product-info">
          <h3>Special Item</h3>
          <div class="rating">
            <span>★★★★★</span>
          </div>
          <router-link to="/restaurante">
            <button class="add-to-cart">Ver Página</button>
          </router-link>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  </script>
  
  <style>
  .product-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 20px;
    justify-content: center;
    padding: 20px;
  }
  
  .product-card {
    border: 1px solid #ddd;
    border-radius: 5px;
    overflow: hidden;
    width: 100%;
    text-align: center;
    font-family: Arial, sans-serif;
    position: relative;
    background: #fff;
  }
  
  .image-placeholder {
    width: 100%;
    height: 200px;
    background: #e1e1e1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #aaa;
    font-size: 20px;
  }
  
  .sale-badge {
    position: absolute;
    top: 10px;
    right: 10px;
    background: #000;
    color: #fff;
    padding: 5px 10px;
    font-size: 12px;
    border-radius: 3px;
  }
  
  .product-info {
    padding: 20px;
  }
  
  .product-info h3 {
    margin: 0;
    font-size: 20px;
  }
  
  .rating span {
    color: gold;
    font-size: 16px;
  }
  
  .add-to-cart {
    background: #007bff;
    color: #fff;
    border: none;
    padding: 10px 20px;
    cursor: pointer;
    border-radius: 5px;
    margin-top: 10px;
  }
  
  .add-to-cart:hover {
    background: #0056b3;
  }
  </style>
  