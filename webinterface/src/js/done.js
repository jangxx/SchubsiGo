const Vue = require("vue/dist/vue.common");
const axios = require("axios");

var app = new Vue({
    data: {
        userinfo: {
            username: "",
            devicename: ""
        },
        logged_in: false,
        registered: false
    },
    methods: {
        logout: function() {
            axios.post("/api/logout").then(resp => {
                window.location = "/"; //redirect to index, so the server can make a descision of where to go next
            }).catch(err => {
                console.log(err);
            });
        }
    },
    created: function() {
        axios.get("/api/userinfo").then(resp => {
            this.userinfo.username = resp.data.username;
            this.userinfo.devicename = resp.data.devicename;
            this.logged_in = resp.data.loggedin;
            this.registered = resp.data.registered;
        });
    }
});

window.addEventListener("DOMContentLoaded", function() {
    app.$mount("main");
});