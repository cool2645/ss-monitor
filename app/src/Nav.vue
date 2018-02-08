<template>
    <div>
        <!-- =============================================== -->

        <!-- Left side column. contains the sidebar for desktop -->
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
        <!-- Left side column. contains the sidebar for mobile -->
        <aside id="menu">
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
    </div>
</template>

<script>
    import config from './config';
    import Slideout from 'slideout';

    export default {
        data() {
            return {
                slideout: null,
                isOpen: false,
            }
        },
        props: ['routes'],
        methods: {
            changeTab(path) {
                this.$router.push(path);
                $('body').removeClass('sidebar-open');
                this.slideout.close();
            }
        },
        mounted() {
            let vm = this;
            document.title = config.appName;
            document.getElementById('logo').innerText = config.appName;
            vm.slideout = new Slideout({
                'panel': document.getElementById('main'),
                'menu': document.getElementById('menu'),
                'padding': 200,
                'tolerance': 70,
            });
            addEventListener('resize', () => {
                vm.slideout.close();
            });
            vm.slideout.on('open', () => {
                vm.isOpen = true;
            });
            vm.slideout.on('close', () => {
                vm.isOpen = false;
            });
            document.getElementById('sidebar-toggle-mobile').addEventListener('click', () => {
                if (vm.isOpen) {
                    vm.slideout.close();
                }
                else {
                    vm.slideout.open();
                }
            });
        }
    }
</script>

<style lang="scss">

</style>
