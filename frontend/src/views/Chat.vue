<template>
  <div class="chat-container">
    <header class="chat-header">
      <a @click="$router.go(-2)" class="back-link">
        <i class="bi bi-chevron-left"></i>
        <span>返回</span>
      </a>
      <h2>{{ query.friendName }}</h2>
    </header>

    <div class="chat-messages" ref="chatMessages">
      <div
        v-for="(msg, index) in messages"
        :key="index"
        :class="['message-container', msg.isOwnMessage ? 'message-own' : 'message-other']"
      >
        <div class="message-content">
          <div v-if="msg.isImage">
            <img :src="msg.content" alt="Image" class="chat-image" />
          </div>
          <div v-else>
            {{ msg.content }}
          </div>
        </div>
        <div :class="['message-time', msg.isOwnMessage ? 'time-left' : 'time-right']">
          {{ formatTime(msg.time) }}
        </div>
      </div>
    </div>

    <div class="chat-input">
      <input type="text" v-model="message" placeholder="Type a message..." @keyup.enter="sendMessage"/>
      <button @click="sendMessage">Send</button>
      <input type="file" @change="uploadImage" accept="image/*" style="display: none;" />
      <button @click="triggerFileUpload">Upload Image</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';

const message = ref('');
const messages = ref([]);
const chatMessages = ref(null); // 引用聊天訊息區域
const { query } = useRoute();
let socket;

// 滾動到聊天區域的底部
const scrollToBottom = () => {
  if (chatMessages.value) {
    chatMessages.value.scrollTop = chatMessages.value.scrollHeight;
  }
};

const sendMessage = () => {
  if (message.value.trim()) {
    // 獲取當前時間並格式化為 YYYY-MM-DD HH:mm
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0'); // 月份是0-11，因此+1
    const day = String(now.getDate()).padStart(2, '0');
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');

    const formattedTime = `${year}-${month}-${day} ${hours}:${minutes}`;

    // 將自己發送的訊息加入訊息列表並靠右顯示
    messages.value.push({ 
      content: message.value, 
      isOwnMessage: true,
      time: formattedTime,
    });

    // 將物件轉換為 JSON 字串並發送給服務器
    socket.send(JSON.stringify({
      TargetId: parseInt(query.friendId),
      Type: 1,
      FromId: parseInt(query.userId),
      Media: 1,
      Content: message.value,
      IsRead: 0, // 0：未讀  1：已讀
      CreateTime: formattedTime,
    }));

    message.value = ''; 
    // 確保 DOM 更新後滾動到最下方
    nextTick(scrollToBottom);
  }
};


const triggerFileUpload = () => {
  // 觸發隱藏的 file input
  document.querySelector('input[type="file"]').click();
};

const uploadImage = async (event) => {
  const file = event.target.files[0];
  if (file) {
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0'); // 月份是0-11，因此+1
    const day = String(now.getDate()).padStart(2, '0');
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');

    const formattedTime = `${year}-${month}-${day} ${hours}:${minutes}`;
    const formData = new FormData(); 
    formData.append('file', file); 

    try {
      const { data } = await axios.post('/user/upload', formData);
      if (data.Code === 0) {
        const imageUrl = data.Data;

        messages.value.push({ 
          content: imageUrl, 
          isOwnMessage: true, 
          isImage: true,
          time: formattedTime,
         });

        socket.send(JSON.stringify({
          TargetId: parseInt(query.friendId),
          Type: 1,
          FromId: parseInt(query.userId),
          Media: 2, 
          Content: imageUrl,
          IsRead: 0, // 0：未讀  1：已讀
          CreateTime: formattedTime,
        }));
        // 確保 DOM 更新後滾動到最下方
        nextTick(scrollToBottom);
      }
    } catch (error) {
      console.error('圖片上傳失敗:', error);
    }
  }
};

// 解析並處理聊天歷史記錄
const parseChatHistory = (history) => {
  history.forEach((item) => {
    const msg = JSON.parse(item);
    const isOwnMessage = msg.FromId === parseInt(query.userId); // 判斷是否是自己發的訊息

    messages.value.push({
      content: msg.Content,
      isOwnMessage,
      isImage: msg.Media === 2, // 判斷是否是圖片訊息
      time: msg.CreateTime,
      isRead: msg.IsRead,
    });
  });
  // 確保 DOM 更新後滾動到最下方
  nextTick(scrollToBottom);
};

const getChatHistory = async () => {
  const params = new URLSearchParams();
  params.append('userIdA', query.userId);
  params.append('userIdB', query.friendId);
  params.append('start', 0);
  params.append('end', -1); 

  try {
    const { data } = await axios.post('/user/getChatHistory', params);
    if (data.Rows) {
      parseChatHistory(data.Rows);
    }
  } catch (error) {
    console.error('獲取聊天記錄失敗:', error);
  }
};


onMounted(() => {
  getChatHistory(); 

  // 建立 WebSocket 連接
  socket = new WebSocket(`ws://localhost:8081/chat?userId=${query.userId}`);

  socket.onopen = () => {
    console.log('WebSocket 連接成功');
  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    const isOwnMessage = data.FromId === parseInt(query.userId); // 判斷是否是自己發的訊息

    messages.value.push({
      content: data.Content,
      isOwnMessage,
      isImage: data.Media === 2,
      time: data.CreateTime,
      isRead: data.IsRead, // 判斷是否是圖片訊息
    });
    // 確保 DOM 更新後滾動到最下方
    nextTick(scrollToBottom);
  };

  socket.onclose = () => {
    console.log('WebSocket 連接關閉');
  };

  socket.onerror = (error) => {
    console.error('WebSocket 連接發生錯誤:', error);
  };
});

const formatTime = (timeStr) => {
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

onBeforeUnmount(() => {
  if (socket) {
    socket.close();
    console.log('組件卸載 關閉WebSocket連接');
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
  justify-content: space-between;
}

.chat-header {
  position: relative;
  background-color: #f8f9fa;
  padding: 1rem;
  border-bottom: 1px solid #ddd;
  display: flex;
  align-items: center;
}

.back-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: inherit;
  cursor: pointer;
}

.chat-header h2 {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  margin: 0;
}

.chat-messages {
  flex-grow: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.message-container {
  display: flex;
  align-items: flex-end;
  max-width: 60%;
}

.message-own {
  flex-direction: row-reverse; 
  align-self: flex-end;
}

.message-other {
  flex-direction: row; 
  align-self: flex-start;
}

.message-content {
  background-color: #daf1da; 
  padding: 0.5rem;
  border-radius: 5px;
  word-wrap: break-word;
}

.message-other .message-content {
  background-color: #f1f0f0;
}

.chat-image {
  max-width: 200px;
  max-height: 200px;
  object-fit: contain;
  border-radius: 5px;
}

.message-time {
  font-size: 0.75rem;
  color: #888;
  margin: 0 0.5rem; 
}

.time-left {
  margin-left: 0.5rem;
}

.time-right {
  margin-right: 0.5rem;
}

.chat-input {
  display: flex;
  gap: 0.5rem;
  padding: 1rem;
  background-color: #f9f9f9;
  border-top: 1px solid #ccc;
}

.chat-input input {
  flex: 1;
  padding: 0.5rem;
}

.chat-input button {
  padding: 0.5rem 1rem;
}
</style>