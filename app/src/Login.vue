<template>
    <!-- Main content -->
    <div>
        <div class="row">
            <div class="col-md-12">
                <div id="msg-success" class="alert alert-success" style="display: none;">
                    <h4><i class="icon fa fa-info"></i> 登录成功，欢迎回来！</h4>

                    <p id="msg-success-p">{{ loginMessage }}
                        <a href="javascript:;" @click="logout">登出</a>
                    </p>
                </div>
            </div>
        </div>
        <div v-if="!isLogin" class="row">
            <!-- left column -->
            <div class="col-md-6">
                <!-- general form elements -->
                <div class="box box-primary">
                    <div class="box-body">
                        <div class="form-horizontal">
                            <div class="row">
                                <fieldset class="col-sm-12">
                                    <legend>管理员登录</legend>
                                    <div class="form-group">
                                        <label for="username" class="col-sm-3 control-label">用户名</label>

                                        <div class="col-sm-9">
                                            <input class="form-control" id="username" v-model="userForm.username">
                                        </div>
                                    </div>
                                    <div class="form-group">
                                        <label for="password" class="col-sm-3 control-label">密码</label>

                                        <div class="col-sm-9">
                                            <input class="form-control" type="password" id="password" @keypress.enter="login" v-model="userForm.password">
                                        </div>
                                    </div>
                                </fieldset>
                            </div>
                        </div>
                    </div>
                    <!-- /.box-body -->
                    <div class="box-footer">
                        <button @click="login" class="btn btn-primary">登录</button>
                    </div>
                </div>
                <!-- /.box -->
            </div>
        </div>
        <!-- /.row -->
    </div>
</template>

<script>
    import config from './config'
    import urlParam from './buildUrlParam'
    const privilege = {
        "admin": "管理员"
    };
    export default {
        props: [
            "auth"
        ],
        data() {
            return {
                warning: "",
                error: "",
                userForm: {
                    username: "",
                    password: ""
                },
            }
        },
        computed: {
            user() {
                return this.auth.user
            },
            isLogin() {
                return this.auth.isLogin
            },
            loginMessage() {
                return "你好，" + this.user.username + "（" + privilege[this.user.privilege] + "）！";
            }
        },
        mounted() {
            if(this.isLogin) {
                this.onLogin()
            }
        },
        watch: {
            isLogin(newV, oldV) {
                if (newV && !oldV) {
                    this.onLogin()
                } else if(!newV && oldV) {
                    this.onLogout()
                }
            }
        },
        methods: {
            onLogin() {
                this.$emit('dismiss-alert');
                $("#msg-success").hide(10).show(100);
            },
            onLogout() {
                $("#msg-success").hide(10);
            },
            login() {
                let vm = this;
                fetch(config.urlPrefix + '/auth?', {
                    credentials: 'include',
                    method: "POST",
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: urlParam(this.userForm)
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                vm.userForm.password = "";
                                if (res.result) {
                                    this.checkStatus();
                                } else {
                                    this.$emit('show-warning', "登录失败：" + res.msg);
                                }
                            }
                        )
                    })
                    .catch(error => {
                        this.$emit('show-error', "发生错误：" + res.msg);
                    });
            },
            checkStatus() {
                this.$emit('check-auth')
            },
            logout() {
                fetch(config.urlPrefix + '/auth?', {
                    credentials: 'include',
                    method: "DELETE",
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    this.checkStatus();
                                } else {
                                    this.$emit('show-warning', "登出失败：" + res.msg);
                                }
                            }
                        )
                    })
                    .catch(error => {
                        this.$emit('show-error', "发生错误：" + res.msg);
                    });
            }
        }
    }
</script>

<style scoped>

</style>