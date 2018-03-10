import Vue from 'vue'
import VueRouter from 'vue-router'
import Nav from './Nav.vue'
import Footer from './Footer.vue'

const Status = () => import('./Status.vue');
const Tasks = () => import('./Tasks.vue');
const Task = () => import('./Task.vue');
const Launch = () => import('./Launch.vue');
const Admin = () => import('./Admin.vue');
const LaravelVuePagination = () => import("laravel-vue-pagination");

import './style.css'
import config from './config'
import TreeView from 'vue-json-tree-view'
import registerServiceWorker from './registerServiceWorker';

registerServiceWorker();

Vue.use(TreeView);
Vue.use(VueRouter);
Vue.component('pagination', LaravelVuePagination);

const routes = [
    {path: '/', component: Status},
    {title: '状态监控', path: '/status', icon: "fa-dashboard", component: Status},
    {title: '任务记录', path: '/task', icon: "fa-flag-checkered", component: Tasks},
    {path: '/task/:id', component: Task},
    {title: '创建任务', path: '/launch', icon: "fa-eye", component: Launch},
    {title: '管理面板', path: '/admin', icon: "fa-server", component: Admin},
];

const router = new VueRouter({
    mode: 'history',
    routes
});

new Vue({
    router,
    el: '#app',
    data: {
        routes: routes,
        auth: {
            isLogin: false,
            user: {
                username: "",
                privilege: ""
            }
        }
    },
    components: {
        'nav-section': Nav,
        'foot-section': Footer
    },
    mounted() {
        this.checkAuth();
    },
    methods: {
        checkAuth() {
            let vm = this;
            fetch(config.urlPrefix + '/auth?', {
                credentials: 'include'
            })
                .then(res => {
                    res.json().then(
                        res => {
                            if (res.result) {
                                vm.auth.user = res.data;
                                vm.auth.isLogin = true;

                            } else {
                                vm.auth.isLogin = false;
                            }
                        }
                    )
                });
        },
        loadAsyncChunks() {
            Status();
            Tasks();
            Task();
            Admin();
            Launch();
            LaravelVuePagination();
        }
    }
});
