import 'package:app/views/result.dart';
import 'package:app/views/upload.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Anti-Fraud User',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: const HomePage(),
    );
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => HomePageState();

  static HomePageState of(BuildContext context) {
    final state = context.findAncestorStateOfType<HomePageState>();
    if (state == null) {
      throw Exception('HomePageState not found');
    }
    return state;
  }
}

class HomePageState extends State<HomePage> {
  int selectedIndex = 0;
  final List<Widget> _pages = [UploadPage(), ResultPage()];

  void updateSelectedPosition(int selectedPosition) {
    setState(() {
      selectedIndex = selectedPosition;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[selectedIndex],
      bottomNavigationBar: NavigationBar(
        backgroundColor: Theme.of(context).colorScheme.background,
        indicatorColor: Theme.of(context).colorScheme.inversePrimary,
        destinations: const [
          NavigationDestination(
              icon: Icon(Icons.map),
              selectedIcon: Icon(Icons.map),
              label: '资料上传'),
          NavigationDestination(
              icon: Icon(Icons.home),
              selectedIcon: Icon(Icons.home),
              label: '分析结果'),
        ],
        selectedIndex: selectedIndex,
        labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
        onDestinationSelected: updateSelectedPosition,
      ),
    );
  }
}
