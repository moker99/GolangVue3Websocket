<template>
  <div>
    <div class="circle-container">
      <div class="circle" @click="showModal('addFriend')">
        <i class="bi bi-person-plus"></i>
        <p>加入好友</p>
      </div>
      <div class="circle" @click="showModal('joinGroup')">
        <i class="bi bi-person-plus"></i>
        <p>加入群組</p>
      </div>
      <div class="circle" @click="showModal('createGroup')">
        <i class="bi bi-person-plus"></i>
        <p>創建群組</p>
      </div>
    </div>

    <!-- Add Friend Modal -->
    <div v-if="activeModal === 'addFriend'" class="modal">
      <div class="modal-content">
        <span class="close" @click="closeModal">×</span>
        <h4>加入好友</h4>
        <input type="text" v-model="newFriendId" placeholder="請輸入好友暱稱：" />
        <div class="button-container">
          <button type="button" class="btn btn-outline-primary small-btn" @click="addFriends">送出</button>
        </div>
      </div>
    </div>

    <!-- Create Group Modal -->
    <div v-if="activeModal === 'createGroup'" class="modal">
      <div class="modal-content">
        <span class="close" @click="closeModal">×</span>
        <h4>創建群組</h4>
        <input type="text" v-model="groupName" placeholder="請輸入群組名稱：" />
        <input type="text" v-model="groupDescription" placeholder="請輸入群組介紹：" />
        <div class="button-container">
          <button type="button" class="btn btn-outline-primary small-btn" @click="createGroup">送出</button>
        </div>
      </div>
    </div>

    <!-- Join Group Modal -->
    <div v-if="activeModal === 'joinGroup'" class="modal">
      <div class="modal-content">
        <span class="close" @click="closeModal">×</span>
        <h4>加入群組</h4>
        <input type="text" v-model="groupId" placeholder="請輸入群組ID：" />
        <div class="button-container">
          <button type="button" class="btn btn-outline-primary small-btn" @click="joinGroup">送出</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';

const newFriendId = ref('');
const groupName = ref('');
const groupDescription = ref('');
const groupId = ref(''); 
const activeModal = ref(null);
const { query } = useRoute();

const showModal = (modalType) => {
  activeModal.value = modalType;
};

const closeModal = () => {
  activeModal.value = null;
};

const addFriends = async () => {
  const params = new URLSearchParams();
  params.append('userId', query.userId);
  params.append('friendName', newFriendId.value);

  if (!newFriendId.value) {
    alert("Please enter a valid friend name.");
    return;
  }

  try {
    const { data } = await axios.post("/user/addFriends", params);
    newFriendId.value = ''; 
    alert(data.Msg);
    closeModal();
  } catch (error) {
    console.error('Error adding friend:', error);
  }
};

const createGroup = async () => {
  const params = new URLSearchParams();
  params.append('ownerId', query.userId);
  params.append('name', groupName.value);
  params.append('desc', groupDescription.value);

  if (!groupName.value || !groupDescription.value) {
    alert("Please enter both group name and description.");
    return;
  }

  try {
    const { data } = await axios.post("/user/createCommunity", params);
    groupName.value = '';
    groupDescription.value = '';
    alert(data.Msg);
    closeModal();
  } catch (error) {
    console.error('Error creating group:', error);
  }
};

const joinGroup = async () => {
  const params = new URLSearchParams();
  params.append('userId', query.userId);
  params.append('comId', groupId.value);

  if (!groupId.value) {
    alert("Please enter a valid group ID.");
    return;
  }

  try {
    const { data } = await axios.post("/user/joinGroup", params);
    groupId.value = '';
    alert(data.Msg);
    closeModal();
  } catch (error) {
    console.error('Error joining group:', error);
  }
};
</script>

<style scoped>
.circle-container {
  width: 100%;
  min-height: 150px;
  padding: 30px;
  margin: 0 auto;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 20px;
}

.circle {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  background-color: #f0f0f0;
  cursor: pointer;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.circle i {
  font-size: 2rem;
}

.circle p {
  margin: 0.5rem 0 0;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  position: relative;
  width: 300px;
}

.modal-content input {
  width: 100%;
  margin-top: 1rem;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.button-container {
  text-align: center;
  margin-top: 1.5rem;
}

.btn {
  width: 100px;
}

.close {
  position: absolute;
  top: 10px;
  right: 10px;
  cursor: pointer;
  font-size: 1.5rem;
}
</style>
