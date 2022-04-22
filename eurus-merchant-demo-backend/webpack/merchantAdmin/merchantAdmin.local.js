// This file is for webpack config
// Check out https://webpack.js.org/configuration/ for more information
// tsconifg.json will ignore this file since this file is .js
const path = require('path');
const webpack = require('webpack');
const nodeExternals = require('webpack-node-externals');
module.exports = {
  entry: {
    server: path.join(
      __dirname,
      '../../dist/tsc/build/merchantAdmin/server/merchantAdmin/index.js',
    ),
  },
  output: {
    path: path.join(__dirname, '../../dist/webpack/local/merchantAdmin'),
    globalObject: 'this',
    publicPath: '/',
    filename: 'bundle.js',
  },
  target: 'node',
  node: {
    // Need this when working with express, otherwise the build fails
    __dirname: false, // if you don't put this is, __dirname
    __filename: false, // and __filename return blank or /
  },
  externals: [nodeExternals()], // Need this to avoid error when working with Express
  module: {
    rules: [
      {
        // Transpiles ES6-8 into ES5
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
        },
      },
      {
        // Loads the javacript into html template provided.
        // Entry point is set below in HtmlWebPackPlugin in Plugins
        test: /\.html$/,
        use: [{ loader: 'html-loader' }],
      },
    ],
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: JSON.stringify('local'),
      },
    }),
  ],
};
