function lm() {

  var m = science.lin.multiply,
    t = science.lin.transpose,
    inv = science.lin.inverse;

  var _y, _X, _data;

  function model() {
    var y = vectorToMatrix(_data.map(_y));  // dependent variable (n x 1 matrix)
    var X = _data.map(_X);                  // independent variables (n x k matrix)
    var n = X.length;                       // n = # of observations
    var k = X[0].length;                    // k = # of model parameters
    var invXX = inv(m(t(X), X));            // (X'X)^1  (k x k matrix)
    var B = m(invXX, m(t(X), y));           // model parameters (k x 1 matrix)
    var y_pred = m(X, B);                   // in-sample model predictions
    var e = d3.range(n).map(function (i) {   // prediction error
      return [y[i][0] - y_pred[i][0]];
    });

    var u = X.map(function (x, i) {
      return x.map(function (d) { return d * e[i][0]; });
    });
    var whites_scaler = n / (n - k);
    var V = m(invXX, m(m(t(u), u), invXX))
      .map(function (row) {
        return row.map(function (d) { return whites_scaler * d; });
      });

    var se = diag(V).map(function (d) { return [Math.sqrt(d[0])]; });

    var tstat = B.map(function (b, i) { return [b[0] / se[i][0]]; });

    return {
      B: B,
      se: se,
      V: V,
      t: tstat,
      getRegressionInterval: getRegressionInterval
    };

    function getRegressionInterval(d) {
      var x = [_X(d)]; // 1 x k
      var y_pred = m(x, B)[0][0];
      var interval = 1.96 * Math.sqrt(m(m(x, V), t(x))[0][0]);
      return {
        y_lower: y_pred - interval,
        y_pred: y_pred,
        y_upper: y_pred + interval
      };
    }
  }

  model.y = function (_) {
    if (!arguments.length) return _y;
    _y = _;
    return model;
  };

  model.X = function (_) {
    if (!arguments.length) return _X;
    _X = _;
    return model;
  };

  model.data = function (_) {
    if (!arguments.length) return _data;
    _data = _;
    return model;
  };

  return model;

  function vectorToMatrix(_) { return _.map(function (d) { return [d]; }); }
  function matrixToVector(_) { return _.map(function (d) { return d[0]; }); }

  function ones(r, c) {
    return d3.range(r).map(function () {
      return d3.range(c).map(function () { return 1; });
    });
  }

  function diag(x) { return x.map(function (d, j) { return [d[j]]; }); }
}

d3.lineChart = function () {
  var width = 300,
    height = 300,
    xExtent, yExtent,
    x, y,
    xScale = d3.scaleLinear(),
    yScale = d3.scaleLinear(),
    path = d3.line();

  function chart(selection, lineClass, data) {
    xScale
      .domain(xExtent || d3.extent(data, x))
      .range([0, width]);

    yScale
      .domain(yExtent || d3.extent(data, y))
      .range([height, 0]);

    xScale.clamp(true);

    path
      .x(function (d) { return xScale(x(d)); })
      .y(function (d) { return yScale(y(d)); });

    var lineSelectorString = "." + (lineClass.replace(" ", "."));

    var lines = selection.selectAll(lineSelectorString).data([data]);

    lines.enter().append("path")
      .attr("class", lineClass);

    lines
      .transition().duration(333)
      .attr("d", path);

    lines.exit().remove();
  }

  chart.width = function (_) {
    if (!arguments.length) return width;
    width = _;
    return chart;
  };

  chart.height = function (_) {
    if (!arguments.length) return height;
    height = _;
    return chart;
  };

  chart.x = function (_) {
    if (!arguments.length) return x;
    x = _;
    return chart;
  };

  chart.y = function (_) {
    if (!arguments.length) return y;
    y = _;
    return chart;
  };

  chart.xExtent = function (_) {
    if (!arguments.length) return xExtent;
    xExtent = _;
    return chart;
  };

  chart.yExtent = function (_) {
    if (!arguments.length) return yExtent;
    yExtent = _;
    return chart;
  };

  // Getter functions
  chart.xAxis = function () { return d3.axisBottom().scale(xScale).ticks(4) };
  chart.yAxis = function () { return d3.axisLeft().scale(yScale).ticks(4) };
  chart.xScale = function () { return xScale; };
  chart.yScale = function () { return yScale; };

  return chart;
}
