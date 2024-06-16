<template>
  <div>
    <div>
      <p>Upload a file</p>
      <button @click="showFileSelect = !showFileSelect">Select a file</button>
    </div>
    <div v-show="showFileSelect">
      <FileUpload :maxSize="1" accept="png" @file-uploaded="getUploadedData" />
    </div>

    <div v-if="fileSelected">
      Successfully Selected file: {{ file.name }}.{{ file.fileExtention }}
    </div>
  </div>
</template>

<script>
import FileUpload from './components/FileUpload.vue'

export default {
  name: 'App',
  components: {
   FileUpload
  },
  data() {
    return {
      file: {},
      fileSelected: false,
      showFileSelect: false,
    };
  },
  methods: {
    getUploadedData(file) {
      this.fileSelected = true;
      this.showFileSelect = false;
      this.file = file;
    },
  },
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
