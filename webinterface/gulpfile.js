const gulp = require('gulp');
const sass = require('gulp-sass');
const ejs = require('gulp-ejs');
const sourcemaps = require('gulp-sourcemaps');
const gulpif = require('gulp-if');
const browserify = require('browserify');
const source = require('vinyl-source-stream');
const buffer = require('vinyl-buffer');
const uglify = require('gulp-uglify-es').default;
const envify = require('envify/custom');
const argv = require('minimist')(process.argv.slice(2));
const glob = require('glob');

const del = require('del');
const fs = require('fs');
const path = require('path');

const paths = {
    styles_base: 'src/css/',
    styles: '*.scss',
    scripts_base: 'src/js/',
    scripts: '*.js',
    html_base: 'src/html/',
    html: '*.ejs'
};

const use_sourcemaps = !argv.production;
const use_uglify = argv.production;

gulp.task('clean', function() {
    return del(['build']);
});

gulp.task('scripts', function() {
    const files = glob.sync(path.join(paths.scripts_base, paths.scripts));
    
    const tasks = files.map(filepath => {
        return new Promise((resolve, reject) => {
            browserify(filepath, {debug: true, extensions: ['es6']})
                .transform(envify({
                    NODE_ENV: (argv.production) ? 'production' : undefined
                }), {global: true})
                .transform("babelify", {presets: ["env"]})
                .bundle()
                .pipe(source(path.basename(filepath)))
                .pipe(buffer())
                .pipe(gulpif(use_sourcemaps, sourcemaps.init({loadMaps: true})))
                .pipe(gulpif(use_uglify, uglify()))
                .pipe(gulpif(use_sourcemaps, sourcemaps.write()))
                .pipe(gulp.dest('./build/js'))
                .on("end", resolve)
                .on("error", reject);
        });
    });

    return Promise.all(tasks);
});

gulp.task('styles', function() {
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
        .pipe(gulp.dest('./build/html'));
});

gulp.task('watch', function() {
    gulp.watch(path.join(paths.scripts_base, paths.scripts), gulp.series(['scripts']));
    gulp.watch(path.join(paths.scripts_base, 'includes', paths.scripts), gulp.series(['scripts']));
    gulp.watch(path.join(paths.styles_base, paths.styles), gulp.series(['styles']));
    gulp.watch(path.join(paths.styles_base, 'includes', paths.styles), gulp.series(['styles']));
    gulp.watch(path.join(paths.html_base, paths.html), gulp.series(['html']));
    gulp.watch(path.join(paths.html_base, 'includes', paths.html), gulp.series(['html']));

    return Promise.resolve();
});

gulp.task('default', gulp.series(['clean', 'scripts', 'styles', 'html', 'watch']));
gulp.task('build', gulp.series(['clean', 'scripts', 'styles', 'html']));
