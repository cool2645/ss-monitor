<template>
    <!-- =============================================== -->

    <!-- Left side column. contains the sidebar -->
    <aside class="main-sidebar">
        <!-- sidebar: style can be found in sidebar.less -->
        <section class="sidebar">
            <!-- Sidebar user panel -->

            <!-- sidebar menu: : style can be found in sidebar.less -->
            <ul class="sidebar-menu">
                <li v-for="route in routes" v-if="route.title" @click="changeTab(route.path)">
                    <a href="javascript:;">
                        <i :class="'fa ' + route.icon"></i> <span>{{route.title}}</span>
                    </a>
                </li>
            </ul>
        </section>
        <!-- /.sidebar -->
    </aside>
</template>

<script>
    import config from './config'
    export default {
        props: [ 'routes' ],
        methods: {
            changeTab (path) {
                this.$router.push(path);
                $('body').removeClass('sidebar-open');
            }
        },
        mounted() {
            document.title = config.appName;
            document.getElementById('logo').innerText = config.appName
        }
    }
</script>

<style lang="scss">

</style>
