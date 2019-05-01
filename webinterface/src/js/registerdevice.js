const Vue = require("vue/dist/vue.common");
const axios = require('axios');

const app = new Vue({
    data: {
        displayError: ""
    },
    methods: {
        submit: function(evt) {
            this.displayError = "";
            let devicename = evt.target.elements.devicename.value;

            axios.post("/api/register", { devicename }).then(resp => {
                window.location = "/"; // redirect to index, so the server can make a descision of where to go next
            }).catch(err => {
                this.displayError = err.response.data;
            });
        }
    }
});

window.addEventListener("DOMContentLoaded", function() {
    app.$mount("main");
});