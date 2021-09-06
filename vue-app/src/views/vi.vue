<template>
  <v-app id="inspire">
    <v-app-bar
      app
    >

      <v-toolbar-title>CHALLENGE</v-toolbar-title>

      <v-spacer></v-spacer>

      
    </v-app-bar>

    <v-main>
      <v-container>
        <v-row>
          <v-col
            v-for="n in 3"
            :key="n"
            cols="4" 
          >
           <router-link to="/v1/syncData">
            <v-card  v-if="n===1" height="300">
              Sync Data </v-card>
           </router-link>
            <v-card @click.prevent="submit" v-if="n===2" height="300">
                <h1>List of clients</h1> </v-card>
            <v-card v-if="n===3" height="300">Client page </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
  //import {reactive} from 'vue';
  export default {
    
    methods:{
        submit:async()=>{
              let response= await fetch('http://localhost:3030/v1/purchaser',{
                method:'GET'
          });
              let data= await response.json
              console.log(data)
              console.log("hola")

          }
      }
  
  }
</script>