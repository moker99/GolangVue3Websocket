<template>
  <div class="container d-flex justify-content-center align-items-center vh-100">
    <div class="card p-4 shadow-lg" style="max-width: 400px; width: 100%;">
      <h2 class="card-title text-center mb-4">Chat Application</h2>
      <div class="mb-3">
        <input type="text" class="form-control" placeholder="暱稱" v-model="user.name">
      </div>
      <div class="mb-3">
        <input 
          type="password" 
          class="form-control" 
          placeholder="密碼" 
          v-model="user.password" 
          @keyup.enter="login"
        >
      </div>
      <button class="btn btn-primary w-100" @click="login">登入</button>
      <div class="links mt-3 d-flex justify-content-between">
        <router-link to="/register" class="text-decoration-none">註冊會員</router-link>
        <router-link to="#" class="text-decoration-none">忘記密碼</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue';
import { useRouter } from 'vue-router'
import axios from 'axios';

const user = reactive({
  name: '',
  password: '',
});

const router = useRouter()

const login = async () => {
  try {
    const params = new URLSearchParams();
    params.append('name', user.name);
    params.append('password', user.password);

    const { data: res } = await axios.post("/user/findUserByNameAndPwd", params);

    if (res.code != 0) {
      alert(res.message);
    } else {
      alert("登入成功,即將跳轉");
      router.push({ name: "Friends", query:{userId:res.data.ID, token:res.data.Identity} });
    }
  } catch (err) {
    console.error("Login failed:", err);
    alert("登入失敗，請稍後重試");
  }
};
</script>

<style></style>
