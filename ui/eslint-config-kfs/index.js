module.exports = {
    "env": {
        "browser": true,
        "es6": true
    },
    "extends": [
        "airbnb",
        "plugin:react-hooks/recommended"
    ],
    "globals": {
        "Atomics": "readonly",
        "SharedArrayBuffer": "readonly"
    },
    "parser": "@babel/eslint-parser",
    "parserOptions": {
        "requireConfigFile": false,
        "babelOptions": {
            "presets": ["@babel/preset-react"]
        },
        "ecmaFeatures": {
            "jsx": true
        },
        "ecmaVersion": 8,
        "sourceType": "module"
    },
    "plugins": [
        "react"
    ],
    "rules": {
        "no-use-before-define": [
            "error",
            {
                "functions": false,
                "classes": false
            }
        ],
        "no-underscore-dangle": "warn",
        "max-len": "warn",
        "no-restricted-syntax": "warn",
        "arrow-parens": "off",
        "no-shadow": "warn",
        "no-param-reassign": "warn",
        "no-unused-expressions": "off",
        "no-unused-vars": "warn",
        "no-bitwise": "off",
        "no-plusplus": "off",
        "no-return-assign": "warn",
        "func-names": "off",
        "arrow-body-style": "warn",
        "prefer-arrow-callback": "warn",
        "prefer-promise-reject-errors": "warn",
        "no-await-in-loop": "off",
        "object-curly-newline": "off",
        "max-classes-per-file": "off",
        "object-property-newline": "warn",
        "import/prefer-default-export": "off",
        "import/no-anonymous-default-export": "off",
        "import/no-mutable-exports": "warn",
        "import/no-extraneous-dependencies": "warn",
        "react/react-in-jsx-scope": "off",
        "react/prop-types": "warn",
        "react/destructuring-assignment": "off",
        "react/state-in-constructor": "off",
        "react/prefer-stateless-function": "off",
        "react/no-array-index-key": "warn",
        "react/static-property-placement": "warn",
        "react/jsx-one-expression-per-line": "warn",
        "react/jsx-filename-extension": "warn",
        "react/jsx-props-no-spreading": "warn",
        "jsx-a11y/no-noninteractive-element-interactions": "warn",
        "jsx-a11y/click-events-have-key-events": "warn"
    },
    "settings": {
        "import/resolver": {
            "node": {
                "paths": [
                    "src"
                ]
            }
        }
    }
}