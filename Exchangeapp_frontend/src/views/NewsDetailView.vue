<template>
  <div class="news-detail-container">
    <div class="content-wrapper">
      <el-card v-if="loading" shadow="never">
        <el-skeleton :rows="6" animated />
      </el-card>

      <el-card v-else-if="article" class="article-detail" shadow="hover">
        <el-button type="primary" link :icon="ArrowLeft" @click="goBack">返回列表</el-button>
        <h1 class="article-title">{{ article.Title }}</h1>
        <p class="article-content">{{ article.Content }}</p>
        <div class="like-row">
          <el-button type="primary" @click="likeArticle" :loading="liking">点赞</el-button>
          <span class="like-count">点赞数：{{ likes }}</span>
        </div>
      </el-card>

      <el-empty v-else description="文章不存在或加载失败" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ArrowLeft } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import axios from '../axios';
import { ApiCode } from '../constants/apiCode';
import { getApiErrorMessage } from '../utils/apiError';
import type { Article } from '../types/Article';

const article = ref<Article | null>(null);
const loading = ref(false);
const liking = ref(false);
const likes = ref(0);
const route = useRoute();
const router = useRouter();

const articleId = String(route.params.id ?? '');

const fetchArticle = async () => {
  if (!articleId) return;
  loading.value = true;
  try {
    const response = await axios.get(`/articles/${articleId}`);
    const payload = response.data;
    if (payload?.code === ApiCode.ok && payload.data) {
      article.value = payload.data as Article;
    } else {
      article.value = null;
    }
  } catch (err) {
    article.value = null;
    ElMessage.error(getApiErrorMessage(err, '加载文章失败'));
  } finally {
    loading.value = false;
  }
};

const fetchLike = async () => {
  if (!articleId) return;
  try {
    const res = await axios.get(`/articles/${articleId}/like`);
    const data = res.data?.data as { likes?: string | number } | undefined;
    likes.value = Number(data?.likes) || 0;
  } catch {
    likes.value = 0;
  }
};

const likeArticle = async () => {
  if (!articleId) return;
  liking.value = true;
  try {
    await axios.post(`/articles/${articleId}/like`);
    await fetchLike();
    ElMessage.success('点赞成功');
  } catch (err) {
    ElMessage.error(getApiErrorMessage(err, '点赞失败'));
  } finally {
    liking.value = false;
  }
};

const goBack = () => {
  router.push({ name: 'News' });
};

onMounted(async () => {
  await fetchArticle();
  await fetchLike();
});
</script>

<style scoped>
.news-detail-container {
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

.article-title {
  margin: 16px 0;
  font-size: 24px;
  color: #303133;
}

.article-content {
  margin: 0 0 24px;
  color: #606266;
  line-height: 1.8;
  white-space: pre-wrap;
}

.like-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.like-count {
  color: #909399;
}
</style>
