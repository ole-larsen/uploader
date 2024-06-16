<template>
  <div class="file-upload">
    <div class="file-upload__area">
      <div v-if="!file.isUploaded">
        <input type="file" name="file" id="file" @change="handleFileChange($event)" />
        <div v-if="errors.length > 0">
          <div
            class="file-upload__error"
            v-for="(error, index) in errors"
            :key="index"
          >
            <span>{{ error }}</span>
          </div>
        </div>
      </div>
      {{ file }}
      <div v-if="file.isUploaded" class="upload-preview">
        <img :src="file.url" v-if="file.isImage" class="file-image" alt="" />
        <div v-if="!file.isImage" class="file-extention">
          {{ file.fileExtention }}
        </div>
        <span class="file-name">
          {{ file.name }}{{ file.isImage ? `.${file.fileExtention}` : "" }}
        </span>
        <div class="">
          <button @click="resetFileInput">Change file</button>
        </div>
        <div class="" style="margin-top: 10px">
          <button @click="sendDataToParent">Select File</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: 'FileUpload',
  props: {
    maxSize: {
      type: Number,
      default: 5,
      required: true,
    },
    accept: {
      type: String,
      default: "image/*",
    },
  },
  async mounted() {
    try {
      await this.load();
    } catch(e) {
      console.error(e);
    }
  },
  data () {
    return {
      errors: [],
      isLoading: false,
      uploadReady: true,
      file: {
        name: "",
        size: 0,
        type: "",
        fileExtention: "",
        url: "",
        isImage: false,
        isUploaded: false,
      },
    };
  },
  methods: {
    isFileSizeValid(fileSize) {
      if (fileSize >= this.maxSize) {
        this.errors.push(`File size should be less than ${this.maxSize} MB`);
      }
    },
    isFileTypeValid(fileExtention) {
      if (!this.accept.split(",").includes(fileExtention)) {
        this.errors.push(`File type should be ${this.accept}`);
      }
    },
    isFileValid(file) {
      this.isFileSizeValid(Math.round((file.size / 1024 / 1024) * 100) / 100);
      //this.isFileTypeValid(file.name.split(".").pop());
      if (this.errors.length === 0) {
        return true;
      } else {
        return false;
      }
    },
    handleFileChange(e) {
      this.errors = [];
      // Check if file is selected
      if (e.target.files && e.target.files[0]) {
       
        // Check if file is valid
        if (this.isFileValid(e.target.files[0])) {
          // Get uploaded file
          const file = e.target.files[0];
            
          this.upload(file)
          
        } else {
          console.log("Invalid file");
        }
      }
    },
    resetFileInput() {
      this.uploadReady = false;
      this.$nextTick(() => {
        this.uploadReady = true;
        this.file = {
          name: "",
          size: 0,
          type: "",
          data: "",
          fileExtention: "",
          url: "",
          isImage: false,
          isUploaded: false,
        };
      });
    },
    sendDataToParent() {
      this.resetFileInput();
      this.$emit("file-uploaded", this.file);
    },
    load () {
      const token = process.env.VUE_APP_XTOKEN;
      console.log(token);
      console.log(process.env);
      return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/v1/uploader/files`, {
        method: 'GET', // or 'PUT'
        headers: {
          "Content-Type": "application/json",
          "x-token": `${token}`
        },
      })
      .then((response) => response.status !== 200 ? response.json() : response.text())
      .then((response) => {
        console.log(response)
      })
      .catch((e) => {
        console.error(e);
        throw e;
      });
    },
    update(file) {
      const token = import.meta.env.VUE_APP_XTOKEN;
      const formData = new FormData();
      formData.append("id", file.id);
      formData.append("name", file.name);
      formData.append("alt", file.alt);
      formData.append("hash", file.hash);
      formData.append("ext", file.ext);
      formData.append("caption", file.caption);
      formData.append("type", file.type);
      formData.append("size", file.size);
      formData.append("width", file.width);
      formData.append("height", file.height);
      formData.append("provider", file.provider);
      formData.append("file", file.file);
      let url = `${import.meta.env.VUE_APP_BACKEND_URL}/api/v1/uploader/files`;

      return fetch(url, {
        method: 'PUT',
        headers: {
          "Authorization": `Bearer ${token}`
        },
        body: formData
      })
      .then((response) => response.status == 200 ? response.json() : response.text())
      .then(async (response) => {
        console.log(response);
      })
      .catch((e) => {
        console.error(e);
        throw e;
      });
    },
    upload(file) {

      const name = file.name.split(".").shift();
      const hash = btoa(name);
      
      const newFile = {
        name:  name,
        size: file.size,
        type: file.type,
        status: "upload",
        hash,
        file
      };
      const token = process.env.VUE_APP_XTOKEN;
   
      const URL = (window.URL || window.webkitURL);
      if (URL) {
        newFile.blob = URL.createObjectURL(newFile.file);
      }
      // Thumbnails
      newFile.thumb = "";
      if (newFile.blob && newFile.type.substr(0, 6) === "image/") {
        newFile.thumb = newFile.blob;
        newFile.alt = newFile.name;
        console.log(newFile.type);
        if (newFile.type === "image/jpeg") {
          newFile.type = "image/jpg";
        }
        newFile.ext = `.${file.name.split(".").pop()}`;
        newFile.caption = newFile.name;
        newFile.cropper = null;
        newFile.provider = "default";
        newFile.status = "thumbed";
      }
  
      if (newFile.blob && newFile.type === "application/pdf") {
        newFile.thumb = newFile.blob;
        
        newFile.ext = `.${file.name.split(".").pop()}`;
        newFile.caption = newFile.name;
        newFile.cropper = null;
        newFile.provider = "default";
        newFile.status = "thumbed";

        return this.reload(token, newFile);
      }
      const img = new Image();
      img.src = newFile.blob;

      img.onload = () => {
        newFile.width = img.width;
        newFile.height = img.height;
        
        return this.reload(token, newFile);
      }
      img.οnerrοr = (err) => {
        console.error(err)
      }
    },
    findByName(name) {
      const token = process.env.VUE_APP_XTOKEN;
      if (name) {
        const url = `${process.env.VUE_APP_BACKEND_URL}/api/v1/uploader/files?name=${name}`;

        return fetch(url, {
          method: 'GET', // or 'PUT'
          headers: {
            "Content-Type": "application/json",
            "x-token": `${token}`
          },
        })
        .then((response) => response.status == 200 ? response.json() : response.text())
        .then((response) => {
          if (response && response.length === 1) {
            return response[0];
          }
          return response;
        })
        .catch((e) => {
          console.error(e);
          throw e;
        });
      } else {
        return null;
      }
    },
    find(id) {
      const token = process.env.VUE_APP_XTOKEN;
      if (id) {
        return fetch(`${process.env.VUE_APP_BACKEND_URL}/api/v1/uploader/files/${id}`, {
          method: 'GET', // or 'PUT'
          headers: {
            "Content-Type": "application/json",
            "x-token": `${token}`
          },
        })
        .then((response) => response.status !== 200 ? response.json() : response.text())
        .then((response) => {
          console.log(response)
          return JSON.parse(response);
        })
        .catch((e) => {
          console.error(e);
          throw e;
        });
      } else {
        return null;
      }
    },
    reload(token, newFile) {
      const formData = new FormData();
      formData.append("name", newFile.name);
      formData.append("alt", newFile.name);
      formData.append("hash", newFile.hash);
      formData.append("ext", newFile.ext);
      formData.append("caption", newFile.caption);
      formData.append("type", newFile.type);
      formData.append("size", newFile.size);
      formData.append("width", newFile.width ? newFile.width : 0);
      formData.append("height", newFile.height ? newFile.height : 0);
      formData.append("provider", newFile.provider);
      formData.append("file", newFile.file);
      
      let url = `${process.env.VUE_APP_BACKEND_URL}/api/v1/uploader/files`;

      return fetch(url, {
        method: 'POST',
        headers: {
          "x-token": `${token}`
        },
        body: formData
      })
      .then((response) => response.status == 200 ? response.json() : response.text())
      .then((response) => {
        return response;
      })
      .catch((e) => {
        console.error(e);
        throw e;
      });
    }
  },
}
</script>


<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.file-upload .file-upload__error {
  margin-top: 10px;
  color: #f00;
  font-size: 12px;
}
.file-upload .upload-preview {
  text-align: center;
}
.file-upload .upload-preview .file-image {
  width: 100%;
  height: 300px;
  object-fit: contain;
}
.file-upload .upload-preview .file-extension {
  height: 100px;
  width: 100px;
  border-radius: 8px;
  background: #ccc;
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 0.5em auto;
  font-size: 1.2em;
  padding: 1em;
  text-transform: uppercase;
  font-weight: 500;
}
.file-upload .upload-preview .file-name {
  font-size: 1.2em;
  font-weight: 500;
  color: #000;
  opacity: 0.5;
}
.file-upload {
  height: 100vh;
  width: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: center;
}
.file-upload .file-upload__area {
  width: 600px;
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px dashed #ccc;
  margin-top: 40px;
}
</style>
