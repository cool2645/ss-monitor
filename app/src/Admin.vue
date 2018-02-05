<template>
    <div>
        <section class="content-header">
            <h1>
                管理面板
                <small>Admin</small>
            </h1>
        </section>
        <section class="content">
            <login :auth="auth" @check-auth="checkAuth"></login>
            <div class="row">
                <node :editing="true"></node>
                <node :editing="false"></node>
            </div>
        </section>
    </div>
</template>

<script>
    import Login from './Login.vue'
    import Node from './Node.vue'
    export default {
        props: [
            "auth"
        ],
        components: {
            "login": Login,
            "node": Node
        },
        methods: {
            checkAuth() {
                this.$emit('check-auth')
            }
        }
    }
</script>

<style scoped>

</style>