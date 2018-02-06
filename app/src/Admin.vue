<template>
    <div>
        <section class="content-header">
            <h1>
                管理面板
                <small>Admin</small>
            </h1>
        </section>
        <section class="content">
            <div class="row">
                <div class="col-md-12">
                    <div id="msg-warning" class="alert alert-warning alert-dismissable" style="display: none;">
                        <button type="button" class="close" @click="dismissAlert" aria-hidden="true">&times;</button>
                        <h4><i class="icon fa fa-warning"></i> 出错了!</h4>

                        <p id="msg-warning-p">{{ warning }}</p>
                    </div>
                    <div id="msg-error" class="alert alert-danger alert-dismissable" style="display: none;">
                        <button type="button" class="close" @click="dismissAlert" aria-hidden="true">&times;</button>
                        <h4><i class="icon fa fa-warning"></i> 出错了!</h4>

                        <p id="msg-error-p">{{ error }}</p>
                    </div>
                </div>
            </div>
            <login :auth="auth" @check-auth="checkAuth"></login>
            <div v-if="isAdmin" class="row">
                <node v-for="node in nodes" :node="node"></node>
            </div>
        </section>
    </div>
</template>

<script>
    import Login from './Login.vue'
    import Node from './Node.vue'
    import config from './config'
    export default {
        props: [
            "auth"
        ],
        data() {
            return {
                nodes: [],
                warning: "",
                error: ""
            }
        },
        components: {
            "login": Login,
            "node": Node
        },
        computed: {
            isAdmin() {
                return this.auth.isLogin && this.auth.user.privilege === "admin"
            }
        },
        mounted() {
            if(this.isAdmin)
                this.getNodes();
        },
        watch: {
            isAdmin(newV, oldV) {
                if (newV && !oldV) {
                    this.getNodes()
                }
            }
        },
        methods: {
            dismissAlert() {
                $("#msg-warning").hide(10);
                $("#msg-error").hide(10);
            },
            checkAuth() {
                this.$emit('check-auth')
            },
            getNodes() {
                let vm = this;
                fetch(config.urlPrefix + '/node?',  {
                    credentials: 'include',
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.nodes = res.data
                                }
                            }
                        )
                    });
            }
        }
    }
</script>

<style scoped>

</style>