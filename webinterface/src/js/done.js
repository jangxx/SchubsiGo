const Vue = require("vue/dist/vue.common");
const axios = require("axios");

const app = new Vue({
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
            if (!confirm("Are you sure you want to log out? Don't forget to remove the device name from your account if you're planning on using it again, later.")) return;

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