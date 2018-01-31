<template>
    <footer class="main-footer">
        <div class="pull-right hidden-xs">
            Made with Love
        </div>
        <strong>Copyright &copy; {{ yearStr }} <a href="#">{{ appName }}</a> </strong>
        All rights reserved. | Powered by <b>ss-monitor</b> v1.0.0
    </footer>
</template>

<script>
    import config from './config'
    export default {
        data() {
            return {
                yearStr: config.yearStr,
                appName: config.appName,
            }
        },
    }
</script>

<style scoped>

</style>