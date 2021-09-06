
<template>
  <v-row justify="center">
    <v-date-picker v-model="date" :key="date">{{date}}</v-date-picker>
    <div class="text-center">
        <v-btn 
          rounded
          color="primary"
          dark
          @click="send"
        >
        Upload Data
        </v-btn>
  </div>
  </v-row>
  
</template>

<script>
  import axios from 'axios'
  export default {
    name:'SyncClients',
    data () {
      return {
        date:"",
      }
    },
    methods:{
        /*send:()=>{
             console.log(this.date)
             fetch('http://localhost:3030/v1',{
                method:'POST',
                headers:{'Content-Type':'application/json'},
                body:JSON.stringify(this.date)
            });
            //await console.log(this.picker)
        }*/
        send(){
          axios.post('http://localhost:3030/v1',{
              date: this.date
          }).then(response=>{
            console.log(response)
          }).catch(err=>{
            console.log(err)
          })
          this.$router.push('/v1')
        }
        /*send(){
          console.log(this.date)
        }*/
    }
    
  }
</script>
<style scoped>

</style>