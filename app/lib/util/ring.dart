import 'dart:math';
import 'package:flutter/material.dart';

class RingChart extends StatelessWidget {
  final double radius;
  final double strokeWidth;
  final double value;
  final Color ringColor;
  final Color progressColor;
  final String text;

  const RingChart(
      {super.key,
      required this.radius,
      required this.strokeWidth,
      required this.value,
      required this.ringColor,
      required this.progressColor,
      required this.text});

  @override
  Widget build(BuildContext context) {
    return CustomPaint(
      size: Size(radius * 2, radius * 2),
      painter: ResourceUsageRingPainter(
          value: value,
          ringColor: ringColor,
          progressColor: progressColor,
          strokeWidth: strokeWidth,
          text: text),
    );
  }
}

class ResourceUsageRingPainter extends CustomPainter {
  final double value;
  final Color ringColor;
  final Color progressColor;
  final double strokeWidth;
  final String text;

  ResourceUsageRingPainter(
      {required this.value,
      required this.ringColor,
      required this.progressColor,
      required this.strokeWidth,
      required this.text});

  @override
  void paint(Canvas canvas, Size size) {
    final double centerX = size.width / 2;
    final double centerY = size.height / 2;
    final double radius =
        min(size.width / 2, size.height / 2) - (strokeWidth / 2);

    const double startAngle = -pi / 2;
    final double sweepAngle = 2 * pi * value;

    final Paint ringPaint = Paint()
      ..color = ringColor
      ..strokeWidth = strokeWidth
      ..style = PaintingStyle.stroke;

    final Paint progressPaint = Paint()
      ..color = progressColor
      ..strokeWidth = strokeWidth
      ..style = PaintingStyle.stroke
      ..strokeCap = StrokeCap.round;

    canvas.drawCircle(Offset(centerX, centerY), radius, ringPaint);
    canvas.drawArc(
      Rect.fromCircle(center: Offset(centerX, centerY), radius: radius),
      startAngle,
      sweepAngle,
      false,
      progressPaint,
    );

    TextSpan span = TextSpan(
      style: const TextStyle(
        color: Colors.grey,
        fontSize: 12.0,
        fontWeight: FontWeight.bold,
      ),
      text: text,
    );
    TextPainter tp = TextPainter(
      text: span,
      textAlign: TextAlign.center,
      textDirection: TextDirection.ltr,
    );
    tp.layout();
    tp.paint(canvas, Offset(centerX - tp.width / 2, centerY - tp.height / 2));
  }

  @override
  bool shouldRepaint(ResourceUsageRingPainter oldDelegate) {
    return oldDelegate.value != value ||
        oldDelegate.ringColor != ringColor ||
        oldDelegate.progressColor != progressColor ||
        oldDelegate.strokeWidth != strokeWidth ||
        oldDelegate.text != text;
  }
}
