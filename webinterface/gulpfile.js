var gulp = require('gulp');
var sass = require('gulp-sass');
var ejs = require('gulp-ejs');
var sourcemaps = require('gulp-sourcemaps');
var gulpif = require('gulp-if');
var browserify = require('browserify');
var source = require('vinyl-source-stream');
var buffer = require('vinyl-buffer');
var uglify = require('gulp-uglify');
var envify = require('envify/custom');
var argv = require('minimist')(process.argv.slice(2));
var glob = require('glob');

var del = require('del');
var fs = require('fs');
var path = require('path');

var paths = {
    styles_base: 'src/css/',
    styles: '*.scss',
    scripts_base: 'src/js/',
    scripts: '*.js',
    html_base: 'src/html/',
    html: '*.ejs'
};

var use_sourcemaps = !argv.production;
var use_uglify = argv.production;

gulp.task('clean', function() {
    return del.sync(['build']);
});

gulp.task('scripts', [], function() {
    var files = glob.sync(path.join(paths.scripts_base, paths.scripts));
    
    var tasks = files.map(function(filepath) {
        return browserify(filepath, {debug: true, extensions: ['es6']})
            .transform(envify({
                NODE_ENV: (argv.production) ? 'production' : undefined
            }), {global: true})
            .transform("babelify", {presets: ["env"]})
            .bundle()
            .pipe(source(path.basename(filepath)))
            .pipe(buffer())
            .pipe(gulpif(use_sourcemaps, sourcemaps.init({loadMaps: true})))
            .pipe(gulpif(use_uglify, uglify().on('error', console.log)))
            .pipe(gulpif(use_sourcemaps, sourcemaps.write()))
            .pipe(gulp.dest('./build/js'));
    });

    return tasks;
});

gulp.task('styles', [], function() {
    return gulp.src(path.join(paths.styles_base, paths.styles))
        .pipe(gulpif(use_sourcemaps, sourcemaps.init()))
            .pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
        .pipe(gulpif(use_sourcemaps, sourcemaps.write()))
        .pipe(gulp.dest('./build/css'));
});

gulp.task('html', function() {
    return gulp.src(path.join(paths.html_base, paths.html))
        .pipe(ejs({

        }, {}, {ext: '.html'}).on('error', console.log))
        .pipe(gulp.dest('./build/html'))
});

gulp.task('watch', function() {
    gulp.watch(path.join(paths.scripts_base, paths.scripts), ['scripts']);
    gulp.watch(path.join(paths.scripts_base, 'includes', paths.scripts), ['scripts']);
    gulp.watch(path.join(paths.styles_base, paths.styles), ['styles']);
    gulp.watch(path.join(paths.styles_base, 'includes', paths.styles), ['styles']);
    gulp.watch(path.join(paths.html_base, paths.html), ['html']);
    gulp.watch(path.join(paths.html_base, 'includes', paths.html), ['html']);
});

gulp.task('default', ['clean', 'scripts', 'styles', 'html', 'watch']);
gulp.task('build', ['clean', 'scripts', 'styles', 'html']);

function getFolders(dir) {
    return fs.readdirSync(dir)
        .filter(function(file) {
            return fs.statSync(path.join(dir, file)).isDirectory();
        });
}

