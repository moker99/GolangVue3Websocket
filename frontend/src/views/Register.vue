<template>
  <div class="container d-flex justify-content-center align-items-center vh-100">
    <div class="card p-4 shadow-lg" style="max-width: 400px; width: 100%;">
      <h2 class="card-title text-center mb-4">註冊會員</h2>
      <div class="mb-3">
        <input type="text" class="form-control" placeholder="暱稱" v-model="user.name">
      </div>
      <div class="mb-3">
        <input type="password" class="form-control" placeholder="密碼" v-model="user.password">
      </div>
      <div class="mb-3">
        <input 
          type="password" 
          class="form-control" 
          placeholder="請再次輸入密碼" 
          v-model="user.repassword"
          @keyup.enter="register"  
        >
      </div>
      <button class="btn btn-primary w-100" @click="register">完成</button>
      <div class="links mt-3 d-flex justify-content-between">
        <router-link to="/" class="text-decoration-none">直接登入</router-link>
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
  repassword: "",
});

const router = useRouter()

const register = async () => {
  try {
    const params = new URLSearchParams();
    params.append('name', user.name);
    params.append('password', user.password);
    params.append('repassword', user.repassword);

    const { data: res } = await axios.post("/user/createUser", params);

    if (res.code != 0) {
      alert(res.message);
    } else {
      alert("註冊成功,即將跳轉");
      router.push({ name: "Friends", query:{userId:res.data.ID, token:res.data.Identity} });
    }
  } catch (err) {
    console.error("Register failed:", err);
    alert("註冊失敗，請稍後重試");
  }
};
</script>