
const app = new Vue({
    el: '#app',
    data: {
        errors: [],
        name: null,
        email: null,
        movie: null,
        passF: null,
        passT: null,
        logreg: true
    },
    methods: {

        checkForm: function (e) {
            this.errors = [];

            if (!this.name) {
                this.errors.push('Укажите имя.');
            }
            if (!(this.passF !== this.passT) || !(this.passF.length == 0)) {
                this.errors.push('неправильный пароль');
            }
            if (this.passF.length == 0) {
                this.errors.push('не введен пароль');
            }
            if (!this.email) {
                this.errors.push('Укажите электронную почту.');
            } else if (!this.validEmail(this.email)) {
                this.errors.push('Укажите корректный адрес электронной почты.');
            }

           
        },
        validEmail: function (email) {
            var re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return re.test(email);
        }
    }

})