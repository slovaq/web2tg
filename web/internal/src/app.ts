import Vue from 'vue'
import VueCookies from 'vue-cookies' // меня заебал линтер, хули он орёт? надо потом вью с аксиосом притащить к себе
import axios from 'axios'
Vue.use(VueCookies)
const app = new Vue({

    delimiters: ['${', '}'],
    el: '#app',
    data: {
        errors: [],
        login: undefined,
        email: null,
        movie: null,
        passF: null,
        passLogin: null,
        passT: null,
        logreg: true,
        RegSuccess:true,
        SuccFalse:true,
        SuccFalseMessage:""
    },
    methods: {
        removeLog(){
            Vue.$cookies.remove("login") 
            Vue.$cookies.remove("password")
            location.href = '/reg'
        },
        axiosLog() {
            console.log("send")
            console.log("this.login: " + this.login + " this.passF: " + this.passF )
            axios.get(`/api/user_get?login=${this.login}&password=${this.passF}`)
                .then(function (response) {
                    // handle success
                    console.log(response);
                    if (response.data.Error!=null){
                        console.log(response.data.Error);
                        Vue.$cookies.remove("login") 
                        Vue.$cookies.remove("password")
                        app.SuccFalse=false
                        app.SuccFalseMessage=response.data.Error
                    } else{
                        console.log("this.login: " + app.login + " this.passF: " + app.passF )
                       // $cookies.set("login", app.login,{ expires: "30d" } );
                       // $cookies.set("password", app.passF,{ expires: "30d" } );
                        console.log("this.login set: " + app.login + " this.passF: " + app.passF)
                        document.cookie = "login=" +app.login+ ";path=/";
                        document.cookie = "password=" +app.passF+ ";path=/";
                    }
                    
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                });

        },
        axiosReg() {
            console.log("send")
            axios.get(`/api/user_create?login=${this.email}&name=${this.login}&password=${this.passF}`)
                .then(function (response) {
                    // handle success
                    console.log(response);
                    document.cookie = `login=${app.email}; expires=30d`;
                    document.cookie = `password=${app.passF}; expires=30d`;
                    // $cookies.set("login", app.email,{ expires: "30d" } );
                    // $cookies.set("password", app.passF,{ expires: "30d" } );
                    
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                });

        },
        logHandler: {
            this.errors = []
            console.log("this.login: ", this.login," \t this.passF: ", this.passF )
            if (!this.login) {
                this.errors.push('Укажите имя.');
            }
           
            if (this.passF.length == 0) {
                this.errors.push('укажите пароль');
            }

            if (this.errors.length == 0) {
                console.log("log ok")
                this.RegSuccess=false
                this.axiosLog()
            }

        },
        subHandler: {
            this.errors = [];
            console.log("this.login: " + this.login + " this.passF: " + this.passF + " this.passT: " + this.passT)
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
            } else if (!this.validEmail(this.email)) {
                this.errors.push('Укажите корректный адрес электронной почты.');
            }
            if (this.errors.length == 0) {
                console.log("reg ok")
                this.RegSuccess=false
                this.axiosReg()
            }

        },
        validEmail: function (email) {
            var re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return re.test(email);
        }
    }

})
