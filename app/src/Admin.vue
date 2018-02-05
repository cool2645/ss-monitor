<template>
    <div>
        <section class="content-header">
            <h1>
                管理面板
                <small>Admin</small>
            </h1>
        </section>
        <section class="content">
            <auth></auth>
            <div class="row">
                <node :editing="true"></node>
                <node :editing="false"></node>
            </div>
        </section>
    </div>
</template>

<script>
    import Auth from './Auth.vue'
    import Node from './Node.vue'
    export default {
        components: {
            "auth": Auth,
            "node": Node
        },
    }
</script>

<style scoped>

</style>