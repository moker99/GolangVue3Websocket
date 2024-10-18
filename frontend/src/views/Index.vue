<template>
  <div class="chat-container">
    <!-- Header -->
    <header class="chat-header">
      <h2>Chat Application</h2>
    </header>

    <!-- Main Content Area -->
    <main class="chat-main">
      <router-view></router-view>
    </main>

    <!-- Footer -->
    <footer class="chat-footer">
      <button :class="{ active: activeTab === 'friends' }" @click="switchTab('friends')">Friends</button>
      <button :class="{ active: activeTab === 'member' }" @click="switchTab('member')">Member</button>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';

const activeTab = ref('friends');
const router = useRouter();
const { query, path } = useRoute();

const switchTab = (tab) => {
  if (tab === 'friends') {
    activeTab.value = 'friends';
    router.push({ name: "Friends", query:{userId: query.userId, token: query.token} });
  }
  else if (tab === 'member') {
    activeTab.value = 'member';
    router.push({ name: "Member", query:{userId: query.userId, token: query.token} });
  } 
};

onMounted(() => {
  if (path.includes('/friends')) {
    activeTab.value = 'friends';
  } else if (path.includes('/member')) {
    activeTab.value = 'member';
  }
});
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
}

.chat-header {
  background-color: #f8f9fa;
  padding: 1rem;
  border-bottom: 1px solid #ddd;
  text-align: center;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  overflow-y: auto;
}

.chat-footer {
  display: flex;
  justify-content: space-around;
  background-color: #f8f9fa;
  border-top: 1px solid #ddd;
  padding: 0.5rem;
}

.chat-footer button {
  border: none;
  background: none;
  padding: 0.5rem;
  cursor: pointer;
}

.chat-footer button.active {
  font-weight: bold;
  border-bottom: 2px solid #007bff;
}

.chat-footer button:hover {
  color: #007bff;
}
</style>
