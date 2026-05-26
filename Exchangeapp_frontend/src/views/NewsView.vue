<template>
  <div class="news-container">
    <div class="content-wrapper">
      <el-card v-if="loading" shadow="never" class="state-card">
        <el-skeleton :rows="4" animated />
      </el-card>

      <template v-else-if="articles.length">
        <el-card v-for="article in articles" :key="article.ID" class="article-card" shadow="hover">
          <h2 class="article-title">{{ article.Title }}</h2>
          <p class="article-preview">{{ article.Preview }}</p>
          <el-button type="primary" link @click="viewDetail(article.ID)">阅读更多</el-button>
        </el-card>
      </template>

      <el-empty v-else description="暂无新闻" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import axios from '../axios';
import { useAuthStore } from '../store/auth';
import type { Article } from '../types/Article';

const articles = ref<Article[]>([]);
const loading = ref(false);
const router = useRouter();
const authStore = useAuthStore();

const fetchArticles = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/articles');
    const payload = response.data;
    if (payload?.code === 0 && Array.isArray(payload.data)) {
      articles.value = payload.data;
    } else if (Array.isArray(payload)) {
      articles.value = payload;
    } else {
      articles.value = [];
    }
  } catch {
    articles.value = [];
    ElMessage.error('加载新闻失败');
  } finally {
    loading.value = false;
  }
};

const viewDetail = (id: string | number) => {
  if (!authStore.isAuthenticated) {
    ElMessage.error('请先登录后再查看');
    return;
  }
  router.push({ name: 'NewsDetail', params: { id: String(id) } });
};

onMounted(fetchArticles);
</script>

<style scoped>
.news-container {
  height: 100%;
  padding: 16px;
  box-sizing: border-box;
  background: #f0f2f5;
  overflow-y: auto;
}

.content-wrapper {
  max-width: 960px;
  margin: 0 auto;
}

.state-card,
.article-card {
  margin-bottom: 16px;
}

.article-title {
  margin: 0 0 12px;
  font-size: 20px;
  color: #303133;
}

.article-preview {
  margin: 0 0 12px;
  color: #606266;
  line-height: 1.6;
}
</style>
