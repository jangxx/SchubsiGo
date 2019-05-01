const Vue = require("vue/dist/vue.common");
const request = require("superagent");

var app = new Vue({
    data: {
        displayError: ""
    },
    methods: {
        submit: function(evt) {
            var self = this;
            var devicename = evt.target.elements.devicename.value;

            request.post("/api/register")
                .type("form")
                .send({devicename})
                .end(function(err, resp) {
                    if (err != null) {
                        self.displayError = resp.text;
                    } else {
                        window.location = "/"; //redirect to index, so the server can make a descision of where to go next
                    }
                });
        }
    }
});

window.addEventListener("load", function() {
    app.$mount("main");
});