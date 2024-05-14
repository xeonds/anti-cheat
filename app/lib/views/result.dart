import 'package:app/util/ring.dart';
import 'package:flutter/material.dart';

class ResultPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('分析结果'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            RingChart(
                radius: 50,
                strokeWidth: 5,
                value: 0.83,
                ringColor: Theme.of(context).colorScheme.primaryContainer,
                progressColor: Theme.of(context).colorScheme.primary,
                text: "83%"),
            const SizedBox(height: 20),
            const Text(
              '结果：疑似诈骗人员',
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
