fetch('/me').then(response => response.json()).then(userInfo => {
    Luigi.setConfig({
        navigation: {
            nodeAccessibilityResolver: (nodeToCheckPermissionFor) => {
                if (nodeToCheckPermissionFor.requiredScopes) {
                    return nodeToCheckPermissionFor.requiredScopes.every(requiredScope => {
                        return userInfo.assignedScopes.some(scope => scope.includes(requiredScope))
                    });
                }
                return true;
            },
            nodes: [{
                pathSegment: 'home',
                context: {
                    accessToken: userInfo.accessToken
                },
                hideFromNav: true,
                children: [{
                    pathSegment: 'home',
                    label: 'Overview',
                    icon: 'hello-world',
                    viewUrl: './vue.js/home.html'
                },
                {
                    pathSegment: 'view',
                    label: 'View',
                    loadingIndicator: { enabled: true },
                    viewGroup: 'viewGroup',
                    category: {
                        label: 'Regular Access',
                        icon: 'employee',
                        collapsible: true
                    },
                    viewUrl: './vue.js/view.html',
                    requiredScopes: ['.View']

                },
                {
                    pathSegment: 'odata',
                    label: 'OData Service',
                    loadingIndicator: { enabled: false },
                    viewGroup: 'viewGroup',
                    category: 'Regular Access',
                    viewUrl: '/ui/faq/',
                    requiredScopes: ['.Admin']
                },
                {
                    viewUrl: './vue.js/admin.html',
                    pathSegment: 'edit-list',
                    loadingIndicator: { enabled: false },
                    viewGroup: 'adminGroup',
                    category: {
                        label: 'Admin Access',
                        icon: 'key',
                        collapsible: true
                    },
                    requiredScopes: ['.Admin'],
                    label: 'Edit (Vue.js)',
                },
                {
                    viewUrl: './fiori-elements/index.html',
                    pathSegment: 'edit-list-report',
                    loadingIndicator: { enabled: false },
                    viewGroup: 'adminGroup',
                    category: 'Admin Access',
                    requiredScopes: ['.Admin'],
                    label: 'Edit (SAP Fiori elements)',
                    category: "Admin Access"
                },
                {
                    viewUrl: './vue.js/add.html',
                    pathSegment: 'add',
                    loadingIndicator: { enabled: false },
                    viewGroup: 'adminGroup',
                    category: 'Admin Access',
                    requiredScopes: ['.Admin'],
                    label: 'Add FAQ',
                    hideFromNav: true,
                    hideSideNav: true,
                    category: "Admin Access"
                },
                {
                    viewUrl: './vue.js/edit.html',
                    pathSegment: 'edit',
                    loadingIndicator: { enabled: false },
                    viewGroup: 'adminGroup',
                    category: 'Admin Access',
                    requiredScopes: ['.Admin'],
                    label: 'Edit FAQ',
                    hideFromNav: true,
                    hideSideNav: true,
                    category: "Admin Access",
                    children: [{
                        viewUrl: './vue.js/edit.html?questionId=:questionId',
                        pathSegment: `:questionId`,
                        context: { questionId: ':questionId' }
                    }]
                },
                {
                    pathSegment: 'export',
                    label: 'Export CSV',
                    viewGroup: 'adminGroup',
                    category: 'Admin Access',
                    externalLink: {
                        url: '/csv/getCSV'
                    },
                    requiredScopes: ['.Admin']
                },
                {
                    label: 'Fundamental Styles',
                    category: {
                        label: 'Useful Links',
                        icon: 'chain-link',
                        collapsible: true
                    },
                    externalLink: {
                        url: 'https://sap.github.io/fundamental-styles/'
                    }
                }, {
                    label: 'Luigi',
                    category: 'Useful Links',
                    externalLink: {
                        url: 'https://luigi-project.io'
                    }
                }]
            }],
            productSwitcher: {
                items: [{
                    icon: 'https://raw.githubusercontent.com/kyma-project/kyma/master/logo.png',
                    label: 'Kyma on GitHub',
                    externalLink: {
                        url: 'https://github.com/kyma-project',
                        sameWindow: false
                    }
                },
                {
                    icon: 'https://cap.cloud.sap/docs/assets/logos/cap.svg',
                    label: 'SAP Cloud Application Programming Model',
                    externalLink: {
                        url: 'https://cap.cloud.sap/docs/',
                        sameWindow: false
                    }
                },{
                    icon: 'https://developers.sap.com/content/dam/application/shared/icons/dev-m-scp-1-start-developing.svg',
                    label: 'SAP BTP Tutorials',
                    externalLink: {
                        url: 'https://developers.sap.com/tutorial-navigator.html?tag=products:technology-platform/sap-business-technology-platform',
                        sameWindow: false
                    }
                },]
            },
            profile: {
                logout: {
                    label: 'Sign Out',
                    customLogoutFn: () => {
                        window.location.href = '/logout';
                    }
                },
                staticUserInfoFn: () => {
                    return { name: userInfo.name }
                }
            }

        },
        routing: {
            useHashRouting: true
        },
        settings: {
            responsiveNavigation2: window.hideShellbar ? 'semiCollapsible' : 'Fiori3',
            responsiveNavigation: 'semiCollapsible',
            header: {
                title: 'Kyma FAQ Management Suite',
                logo: './sap-logo.png'
            }
        }
    });
});
