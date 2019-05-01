const Vue = require("vue/dist/vue.common");
const axios = require("axios");

const app = new Vue({
    data: {
        displayError: "",
        require2FA: false,
    },
    methods: {
        submit: function(evt) {
            this.displayError = "";

            let email = evt.target.elements.email.value;
            let password = evt.target.elements.password.value;

            let postData = { email, password };

            if (this.require2FA) {
                postData.twofacode = evt.target.elements.twofacode.value;
            }

            axios.post("/api/login", postData).then(resp => {
                if (resp.data == "2FA_MISSING") {
                    this.require2FA = true;
                } else {
                    window.location = "/"; // redirect to index, so the server can make a decision of where to go next
                }
            }).catch(err => {
                this.displayError = err.response.data;
            });
        }
    }
});

window.addEventListener("DOMContentLoaded", function() {
    app.$mount("main");
});