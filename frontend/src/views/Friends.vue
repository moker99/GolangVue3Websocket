<template>
  <div>
    <h3>好友列表</h3>
    <div class="list-group" v-if="friends.length">
      <a href="#" v-for="(friend, index) in friends" :key="index" @click="goToChat(friend)" class="list-group-item list-group-item-action" aria-current="true">
        <div class="d-flex w-100 justify-content-between">
          <h5 class="mb-1">{{ friend.Name }}</h5>
          <small>{{ formatTime(friend.LatestMessageTime) }}</small>
        </div>
        <p class="mb-1">{{ friend.LatestMessage !== '' ? formatMessage(friend.LatestMessage) : '' }}</p>
      </a>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';

const friends = ref([]);
const groups = ref([]);
const router = useRouter();
const { query } = useRoute();
let socket;

const loadFriends = async () => {
  const params = new URLSearchParams();
  params.append('userId', query.userId);

  try {
    const { data } = await axios.post("/user/loadFriends", params);
    if (data.Rows) {
      data.Rows.sort((a, b) => {
        if (a.LatestMessageTime > b.LatestMessageTime) return -1;
        if (a.LatestMessageTime < b.LatestMessageTime) return 1;
        return 0;
      });
    }
    friends.value = data.Rows || [];
  } catch (error) {
    console.error('Error searching friends:', error);
  }
};

const loadGroups = async () => {
  const params = new URLSearchParams();
  params.append('ownerId', query.userId);

  try {
    const { data } = await axios.post("/user/loadCommunity", params);
    groups.value = data.Rows || [];
  } catch (error) {
    console.error('Error searching groups:', error);
  }
};

const goToChat = (friend) => {
  router.push({ name: 'Chat', query: { userId: query.userId, friendId: friend.ID, friendName: friend.Name } });
};

const formatTime = (timeStr) => {
  if(timeStr == '') return;
  const date = new Date(timeStr);
  const hours = date.getHours(); // 24小時制的時
  const minutes = String(date.getMinutes()).padStart(2, '0'); // 分鐘

  let period = ''; // 時段標籤
  let formattedHours = hours % 12 || 12; // 12 小時制，00:00-11:59 對應 12:00-11:59

  // 根據不同時間段設置時段標籤
  if (hours >= 0 && hours < 6) {
    period = '凌晨';
  } else if (hours >= 6 && hours < 12) {
    period = '上午';
  } else if (hours >= 12 && hours < 18) {
    period = '下午';
  } else {
    period = '晚上';
  }

  return `${period} ${formattedHours}:${minutes}`; // 返回12小時制時間加時段
};

const formatMessage = (Str) => {
  const imageExtensions = ['.jpg', '.png', '.jpeg'];
  
  if (imageExtensions.some(ext => Str.toLowerCase().endsWith(ext))) {
    return '[圖片]';
  }

  return Str; 
};

onMounted(() => {
  loadFriends();
  loadGroups();

  // 建立 WebSocket 連接
  socket = new WebSocket(`ws://localhost:8081/chat?userId=${query.userId}`);

  socket.onopen = () => {
    console.log('WebSocket 連接成功');
  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);

    const friend = friends.value.find(friend => friend.ID === data.FromId);

    if (friend) {
      friend.LatestMessage = data.Content || ''; 
      friend.LatestMessageTime = data.CreateTime; 
      
      friends.value.sort((a, b) => {
        if (a.LatestMessageTime > b.LatestMessageTime) return -1;
        if (a.LatestMessageTime < b.LatestMessageTime) return 1;
        return 0;
      });
    }
  };

  socket.onclose = () => {
    console.log('WebSocket 連接關閉');
  };

  socket.onerror = (error) => {
    console.error('WebSocket 連接發生錯誤:', error);
  };
});

onBeforeUnmount(() => {
  if (socket) {
    socket.close();
    console.log('組件卸載 關閉WebSocket連接');
  }
});
</script>

<style scoped>
ul {
  list-style-type: none;
  padding: 0;
}

li {
  padding: 0.5rem 0;
  border-bottom: 1px solid #ccc;
  background-color: yellow;
  cursor: pointer;
}

li:hover {
  background-color: #f0f0f0;
}

.list-group-item{
  height: 70px;
}
</style>
