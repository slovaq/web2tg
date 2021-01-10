
const app = new Vue({
    el: '#app',
    data: {
        errors: [],
        login: null,
        email: null,
        movie: null,
        passF: null,
        passT: null,
        logreg: true
    },
    methods: {

        subHandler() {
            this.errors = [];
console.log("this.login: "+this.login+" this.passF: "+this.passF+" this.passT: "+this.passT)
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
            if(this.errors.length==0){
                console.log("reg ok")
            }
           
        },
        validEmail: function (email) {
            var re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return re.test(email);
        }
    }

})