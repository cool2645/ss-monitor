import Vue from 'vue'
import VueRouter from 'vue-router'
import Nav from './Nav.vue'
import Footer from './Footer.vue'
import Status from './Status.vue'
import Tasks from './Tasks.vue'
import Task from './Task.vue'
import Launch from './Launch.vue'
import Admin from './Admin.vue'
import TreeView from "vue-json-tree-view"
import LaravelVuePagination from 'laravel-vue-pagination'
import './style.css'
import config from './config'
import registerServiceWorker from './registerServiceWorker';
registerServiceWorker();

Vue.use(TreeView);
Vue.use(VueRouter);
Vue.component('pagination', LaravelVuePagination);

const routes = [
    { path: '/', redirect: '/status' },
    { title: '状态监控', path: '/status', icon: "fa-dashboard", component: Status},
    { title: '任务记录', path: '/task', icon: "fa-flag-checkered", component: Tasks },
    { path: '/task/:id', component: Task },
    { title: '创建任务', path: '/launch', icon: "fa-eye", component: Launch },
    { title: '管理面板', path: '/admin', icon: "fa-server", component: Admin },
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
    }
});
