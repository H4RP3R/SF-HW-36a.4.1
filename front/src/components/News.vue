<template>
  <div>
    <h2 class="mx-5 my-5">GoNews - агрегатор новостей.</h2>
    <div v-for="post in news" :key="post.ID">
      <v-card elevation="10" outlined class="mx-5 my-5">
        <v-card-title>
          <a :href="post.link" target="_blank"> {{ post.title }} </a>
        </v-card-title>
        <v-card-text>
          <!-- Render HTML content -->
          <div v-html="post.content" class="news-content"></div>
          <v-card-subtitle>
            {{ new Date(post.published).toLocaleString() }}
          </v-card-subtitle>
        </v-card-text>
      </v-card>
    </div>
  </div>
</template>

<script>
export default {
  name: "News",
  data() {
    return {
      news: [],
    };
  },
  mounted() {
    let url = "http://localhost:8088/news/20";
    fetch(url)
      .then((response) => response.json())
      .then((data) => (this.news = data));
  },
};
</script>

<style scoped></style>