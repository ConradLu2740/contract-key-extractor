<template>
  <div class="app-container">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>Contract Key Information Extractor</h1>
          <p class="subtitle">Legal Contract Analysis System for Law Researchers</p>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
      <el-footer>
        <p>ISAC-CRAN System - Contract Key Information Extraction</p>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup>
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.app-container {
  min-height: 100vh;
}

.el-container {
  min-height: 100vh;
}

.el-header {
  background: rgba(255, 255, 255, 0.95);
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.header-content {
  text-align: center;
}

.header-content h1 {
  color: #333;
  font-size: 28px;
  margin-bottom: 8px;
}

.subtitle {
  color: #666;
  font-size: 14px;
}

.el-main {
  padding: 30px;
  background: transparent;
}

.el-footer {
  background: rgba(255, 255, 255, 0.9);
  text-align: center;
  padding: 15px;
  color: #666;
}
</style>
