"use strict";
exports.__esModule = true;
var axios_1 = require("axios");
var vue_1 = require("vue");
var app = new vue_1["default"]({
    delimiters: ['${', '}'],
    el: '#app',
    data: {
        errors: [],
        login: null,
        email: null,
        movie: null,
        passF: null,
        passLogin: null,
        passT: null,
        logreg: true,
        RegSuccess: true,
        SuccFalse: true,
        SuccFalseMessage: ""
    },
    methods: {
        removeLog: function () {
            $cookies.remove("login");
            $cookies.remove("password");
            location.href = '/reg';
        },
        axiosLog: function () {
            console.log("send");
            console.log("this.login: " + this.login + " this.passF: " + this.passF);
            axios_1["default"].get("/api/user_get?login=" + this.login + "&password=" + this.passF)
                .then(function (response) {
                // handle success
                console.log(response);
                if (response.data.Error != null) {
                    console.log(response.data.Error);
                    $cookies.remove("login");
                    $cookies.remove("password");
                    app.SuccFalse = false;
                    app.SuccFalseMessage = response.data.Error;
                }
                else {
                    console.log("this.login: " + app.login + " this.passF: " + app.passF);
                    // $cookies.set("login", app.login,{ expires: "30d" } );
                    // $cookies.set("password", app.passF,{ expires: "30d" } );
                    console.log("this.login set: " + app.login + " this.passF: " + app.passF);
                    document.cookie = "login=" + app.login + ";path=/";
                    document.cookie = "password=" + app.passF + ";path=/";
                }
            })["catch"](function (error) {
                // handle error
                console.log(error);
            });
        },
        axiosReg: function () {
            console.log("send");
            axios_1["default"].get("/api/user_create?login=" + this.email + "&name=" + this.login + "&password=" + this.passF)
                .then(function (response) {
                // handle success
                console.log(response);
                document.cookie = "login=" + app.email + "; expires=30d";
                document.cookie = "password=" + app.passF + "; expires=30d";
                // $cookies.set("login", app.email,{ expires: "30d" } );
                // $cookies.set("password", app.passF,{ expires: "30d" } );
            })["catch"](function (error) {
                // handle error
                console.log(error);
            });
        },
        logHandler: function () {
            this.errors = [];
            console.log("this.login: " + this.login + " this.passF: " + this.passF);
            if (!this.login) {
                this.errors.push('Укажите имя.');
            }
            if (this.passF.length == 0) {
                this.errors.push('укажите пароль');
            }
            if (this.errors.length == 0) {
                console.log("log ok");
                this.RegSuccess = false;
                this.axiosLog();
            }
        },
        subHandler: function () {
            this.errors = [];
            console.log("this.login: " + this.login + " this.passF: " + this.passF + " this.passT: " + this.passT);
            if (!this.login) {
                this.errors.push('Укажите имя.');
            }
            if (this.passF != this.passT) {
                this.errors.push('пароли не совпадают');
            }
            if (this.passF.length == 0) {
                this.errors.push('укажите пароль');
            }
            if (!this.email) {
                this.errors.push('Укажите электронную почту.');
            }
            else if (!this.validEmail(this.email)) {
                this.errors.push('Укажите корректный адрес электронной почты.');
            }
            if (this.errors.length == 0) {
                console.log("reg ok");
                this.RegSuccess = false;
                this.axiosReg();
            }
        },
        validEmail: function (email) {
            var re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return re.test(email);
        }
    }
});
