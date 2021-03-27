
Vue.component('vue-ctk-date-time-picker', window['vue-ctk-date-time-picker']);
var app = new Vue({
    el: '#app',
    data: {
        date: null,
        city: 'city',
        repeat: "one",
        login: null,
        message: null,
        modVisible: false,
        value: null,
        passT: null,
        HttpError: false,
        locale: 'en',
        items: [],
        picked: false,
        file: '',
        checkedNames: [],
    },
    created() {
        console.log("send");
        axios.get(`/auth/record_get`)
            .then(function (response) {
            console.log(response);
            if (response.data.Error != null) {
                console.log(response.data.Error);
            }
            else {
                console.log("created ok");
                console.log(decodeURIComponent(response.data));
                app.items = response.data;
            }
        })
            .catch(function (error) {
            console.log(error);
        });
    },
    methods: {
        handleFileUpload() {
            if (this.$refs.file.files[0].type !== "image/png" && this.$refs.file.files[0].type !== "image/jpeg"
                && this.$refs.file.files[0].type !== "image/jpg") {
                alert("Картинка должна быть типа png или jpeg/jpg. измените пожалуйста картинку. ");
            }
        },
        inst() {
            app.modVisible = true;
        },
        deletepost(id) {
            console.log("delete>", id);
            axios.get(`/auth/record_delete?id=${id}`)
                .then(function (response) {
                console.log("deleted> N");
                console.log(response);
                if (response.data.Error != null) {
                    console.log(response.data.Error);
                    console.log("deleted> N");
                }
                else {
                    console.log("deleted> ok");
                    app.updatePosts();
                }
            })
                .catch(function (error) {
                console.log("deleted> N");
                console.log(error);
            });
        },
        updatePosts() {
            console.log("updatePosts> send");
            axios.get(`/auth/record_get`)
                .then(function (response) {
                console.log(response);
                if (response.data.Error != null) {
                    console.log(response.data.Error);
                }
                else {
                    console.log("updatePosts> created ok");
                    console.log(decodeURIComponent(response.data));
                    app.items = response.data;
                }
            })
                .catch(function (error) {
                console.log(error);
            });
        },
        AxiosSend(message, date, city, repeat) {
            console.log("AxiosSend> message:" + message + " date: " + date + " city: " + city + " repeat: " + repeat);
            var formData = new FormData();
            let file = this.$refs.file.files[0];
            if (file != undefined) {
                if (file.type !== "image/png" && file.type !== "image/jpeg" && file.type !== "image/jpg") {
                    alert("send: Картинка должна быть типа png или jpeg/jpg. измените пожалуйста картинку. ");
                    return;
                }
                formData.append("file", file);
            }
            if (app.picked == false) {
                repeat = "one";
            }
            else {
                repeat = "week";
            }
            axios.post(`/auth/record_create?message=${message}&date=${date}&city=city&period=${repeat}&week=${app.checkedNames}`, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })
                .then(function (response) {
                console.log(response);
                if (response.data.Error != null) {
                    console.log("error:>", response.data.Error);
                }
                else {
                    if (response.data.HttpError == "link") {
                        console.log("error:>", response.data.HttpError);
                        app.modVisible = false;
                        app.HttpError = true;
                        app.modVisible = true;
                    }
                    else {
                        console.log("AxiosSend> else> ok");
                        app.updatePosts();
                        this.message = "";
                        this.date = null;
                        app.HttpError = false;
                        app.modVisible = false;
                    }
                }
            })
                .catch(function (error) {
                console.log(error);
            });
            app.checkedNames = [];
            app.file = '';
        },
        subHandler() {
            this.errors = [];
            console.log(" message:" + this.message + " date: " + this.date + " city: " + this.city + " repeat: " + this.repeat);
            if (!this.message) {
                this.errors.push('введите сообщение');
            }
            if (this.errors.length == 0) {
                console.log("reg ok");
                app.AxiosSend(this.message, this.date, this.city, this.repeat);
            }
            else {
                console.log(this.errors);
            }
        },
    }
});
//# sourceMappingURL=posts.js.map