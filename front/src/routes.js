import HomePage from './components/pages/HomePage';
import NotFoundPage from './components/pages/NotFoundPage';
import PanelLeftPage from './components/pages/PanelLeftPage';
import LearnPage from './components/pages/LearnPage';
import ReviewPage from './components/pages/ReviewPage';
import LoginPage from './components/pages/LoginPage';
import ReviewSelection from './components/pages/ReviewSelection';


export default [
  {
    path: '/',
    component: HomePage,
  },
  {
    path: '/login/',
    component: LoginPage,
  },
  {
    path: '/learn/',
    component: LearnPage,
  },
  {
    path: '/review/',
    component: ReviewPage,
    options: {
      props: {
        learn: [],
      },
    },
  },
  {
    path: '/reviewselection/',
    component: ReviewSelection,
  },
  {
    path: '/panel-left/',
    component: PanelLeftPage,
  },
  {
    path: '(.*)',
    component: NotFoundPage,
  },
];
